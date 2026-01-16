format: 
	@go mod tidy -e
	@go vet ./...
	@gofmt -s -w .

generate: format
	protoc --go_out=. --go-grpc_out=. ./resources/internal.proto 

build: 
	docker build -t github.com/h2hsecure/netbox:latest .

local: build
	docker run --rm -it --name worker --env-file=local.env -p 8080:80 -p 2112:2112 github.com/h2hsecure/netbox:latest

test: format
	go test -cover ./...

treetest: format
	@rm cover.out out.svg || true
	@go test -coverprofile cover.out ./...
	@go-cover-treemap -coverprofile cover.out > out.svg
	@open -a Safari out.svg

audit: format
	@go tool staticcheck -checks=all,-ST1000,-ST1001,-ST1003,-ST1005,-SA1019,-ST1020,-ST1021,-ST1022 ./...
	@go tool govulncheck ./...