format: 
	@go mod tidy -e
	@go vet ./...
	@gofmt -s -w .

generate: format
	protoc --go_out=. --go-grpc_out=. ./resources/internal.proto 

build: 
	docker build -t register.h2hsecure.com/ddos/worker:latest .

local: build
	docker run --rm -it --name worker --env-file=local.env -p 8080:80 -p 2112:2112 register.h2hsecure.com/ddos/worker:latest
