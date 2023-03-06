# URL Shortener

This is a simple URL shortener written in Go using AWS DynamoDB and Lambda.

## AWS Setup

To start you need to create a AWS DynamoDB table. You can do it using AWS Console or AWS CLI.

### Do it using Terraform

As well, you can use Terraform to create all necessary AWS resources.

First, you need to build the application:

```bash
GOOS=linux GOARCH=amd64 go build -o build/main ./cmd/lambda/main.go
```

You need to update Terraform configuration in [`terraform/main.tf`](terraform/main.tf#L1-L20) with your own S3 bucket name and AWS Region.
Don't forget to create a S3 bucket with name you specified in Terraform configuration.

Now you can apply Terraform configuration:

```bash
cd terraform
terraform init
terraform apply
```

Command `terraform apply` will create:

- DynamoDB table with name `url-shortener-table`;
- API Gateway with name `url-shortener-api`;
- Lambda function with name `url-shortener-lambda`.

## How to start (locally)

Execute next commands to start application:

```bash
go mod download
go run ./cmd/app/main.go
```
