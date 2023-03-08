build:
	go build -o ./build/main ./cmd/lambda/main.go

shorten:
	k6 run -e BASE_URL=${BASE_URL} test/load/shorten_url.js 

geturl:
	k6 run -e BASE_URL=${BASE_URL} test/load/get_url.js 

unittest:
	go test -v ./test/unit/...

integrationtest:
	go test -v ./test/integration/...

e2etest:
	go test -v ./test/e2e/...

test:
	go test -v ./test/...