terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  backend "s3" {
    bucket = "local-dm-tfstate"
    key    = "url-shortener/terraform.tfstate"
    region = "eu-central-1"
  }

  required_version = ">= 1.3.6"
}

provider "aws" {
  region = "eu-central-1"
}

data "aws_iam_policy_document" "lambda_policy" {
  statement {
    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]
    effect    = "Allow"
    resources = ["arn:aws:logs:*:*:*"]
  }

  statement {
    actions = [
      "dynamodb:PutItem",
      "dynamodb:DeleteItem",
    ]
    effect    = "Allow"
    resources = [aws_dynamodb_table.table.arn]
  }
}

data "aws_iam_policy_document" "api_policy" {
  statement {
    actions = [
      "lambda:InvokeFunction",
    ]
    effect    = "Allow"
    resources = [aws_lambda_function.lambda.arn]
  }
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "../build/main"
  output_path = "../lambda.zip"
}

resource "aws_iam_policy" "lambda_policy" {
  name   = "LambdaPolicy"
  path   = "/"
  policy = data.aws_iam_policy_document.lambda_policy.json
}

resource "aws_iam_policy" "api_policy" {
  name   = "APIGatewayPolicy"
  path   = "/"
  policy = data.aws_iam_policy_document.api_policy.json
}

resource "aws_iam_role" "lambda_role" {
  name = "LambdaRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })

  managed_policy_arns = [aws_iam_policy.lambda_policy.arn]
}

resource "aws_iam_role" "api_role" {
  name = "APIGatewayRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "apigateway.amazonaws.com"
        }
      },
    ]
  })

  managed_policy_arns = [aws_iam_policy.api_policy.arn]
}

resource "aws_lambda_function" "lambda" {
  filename         = data.archive_file.lambda_zip.output_path
  function_name    = "url-shortener"
  role             = aws_iam_role.lambda_role.arn
  handler          = "main"
  runtime          = "go1.x"
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  environment {
    variables = {
      BASE_URL = aws_apigatewayv2_stage.stage.invoke_url
      TABLE_NAME = aws_dynamodb_table.table.name
      STAGE_NAME = local.stage_name
    }
  }
}

resource "aws_cloudwatch_log_group" "logs" {
  name              = "/aws/lambda/${aws_lambda_function.lambda.function_name}"
  retention_in_days = 30
}

resource "aws_apigatewayv2_api" "api" {
  name          = "url-shortener-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "integration" {
  api_id           = aws_apigatewayv2_api.api.id
  integration_type = "AWS_PROXY"
  description      = "URL Shortener API"
  integration_uri  = aws_lambda_function.lambda.invoke_arn
  credentials_arn  = aws_iam_role.api_role.arn
}

resource "aws_apigatewayv2_route" "route" {
  api_id    = aws_apigatewayv2_api.api.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.integration.id}"
}

resource "aws_apigatewayv2_stage" "stage" {
  api_id      = aws_apigatewayv2_api.api.id
  name        = local.stage_name
  auto_deploy = true
}

resource "aws_lambda_permission" "lambda_permissions" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}

resource "aws_dynamodb_table" "table" {
  name         = "url-shortener-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "Key"

  attribute {
    name = "Key"
    type = "S"
  }

  ttl {
    attribute_name = "Ttl"
    enabled        = true
  }
}
