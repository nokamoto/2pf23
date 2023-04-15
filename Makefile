PROTO_FILES = $(shell find . -name '*.proto')

all: proto go

go:
	go install mvdan.cc/gofumpt@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	gofumpt -l -w .
	go vet ./...
	staticcheck ./...
	go test ./...
	go mod tidy

proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	go install github.com/bufbuild/buf/cmd/buf@latest
	buf lint --error-format=json | jq .
	buf format -w
	buf generate

cialpha:
	go run ./build/ci/test/main.go

.PHONY: all go proto cialpha
