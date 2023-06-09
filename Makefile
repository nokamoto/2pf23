COMMANDS = pf ke-apis
GOBIN = $(shell go env GOPATH)/bin

define ko
	echo $1
	if [ -z "$$KO_DOCKER_REPO" ]; then \
		KO_DOCKER_REPO=ghcr.io/nokamoto/2pf23 ko build --base-import-paths ./cmd/$1; \
	else \
		ko build --base-import-paths ./cmd/$1; \
	fi

endef

fast:
	go test ./...

test: proto testdata gen go

lint: $(GOBIN)/golangci-lint
	golangci-lint run

all: test lint build

$(GOBIN)/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin

$(GOBIN)/gofumpt:
	go install mvdan.cc/gofumpt@v0.5.0

$(GOBIN)/protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0

$(GOBIN)/protoc-gen-connect-go:
	go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@v1.7.0

$(GOBIN)/buf:
	go install github.com/bufbuild/buf/cmd/buf@v1.17.0

$(GOBIN)/ko:
	go install github.com/google/ko@v0.14.1

go: $(GOBIN)/gofumpt
	go generate ./...
	gofumpt -l -w .
	go test ./... -race -covermode=atomic -coverprofile=coverage.out
	go mod tidy

proto: $(GOBIN)/protoc-gen-go $(GOBIN)/buf $(GOBIN)/protoc-gen-connect-go
	buf lint --config build/buf/buf.yaml --error-format=json
	buf format --config build/buf/buf.yaml -w
	buf generate --template build/buf/buf.gen.yaml

cialpha:
	go run ./build/ci/test/main.go

.PHONY: build
build: $(GOBIN)/ko
	$(foreach command,$(COMMANDS),$(call ko,$(command)))

.PHONY: testdata
testdata:
	go run ./cmd/tools cli testdata/generator/cli/test.json internal/generator/cli/generated github.com/nokamoto/2pf23/internal/generator/cli/generated
	go run ./cmd/tools server testdata/generator/server internal/generator/server/generated --mock
	go run ./cmd/tools ent testdata/generator/ent internal/generator/ent/generated generated

gen: $(GOBIN)/buf
	go install ./cmd/protoc-gen-cli
	go install ./cmd/protoc-gen-server
	go install ./cmd/protoc-gen-ent
	buf generate --template build/buf/buf.gen.local.yaml
	go run ./cmd/tools cli build/cli/test.json internal/cli/generated github.com/nokamoto/2pf23/internal/cli/generated
	go run ./cmd/tools server build/server internal/server/generated
	go run ./cmd/tools ent build/ent internal/ent/proto proto

pack:
	sudo add-apt-repository -y ppa:cncf-buildpacks/pack-cli
	sudo apt-get update
	sudo apt-get install pack-cli

tilt: pack $(GOBIN)/ko
	curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
