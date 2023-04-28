API_GO_FILES = $(shell find pkg/api -name '*.go')
COMMANDS = pf ke-apis

define ko
	echo $1
	if [ -z "$$KO_DOCKER_REPO" ]; then \
		KO_DOCKER_REPO=ghcr.io/nokamoto/2pf23 ko build --base-import-paths ./cmd/$1; \
	else \
		ko build --base-import-paths ./cmd/$1; \
	fi

endef

all: proto mock testdata cli go

go:
	go install mvdan.cc/gofumpt@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	gofumpt -l -w .
	go vet ./...
	staticcheck ./...
	go test ./... -race -covermode=atomic -coverprofile=coverage.out
	go mod tidy

proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	go install github.com/bufbuild/buf/cmd/buf@latest
	buf lint --error-format=json
	buf format -w
	buf generate

cialpha:
	go run ./build/ci/test/main.go

mock:
	go install github.com/golang/mock/mockgen@v1.6.0
	go generate ./...
	$(foreach file,$(API_GO_FILES),mockgen -source $(file) -destination internal/mock/$(file))

build:
	go install github.com/google/ko@latest
	go install ./cmd/pf
	go install ./cmd/cli-gen
	$(foreach command,$(COMMANDS),$(call ko,$(command)))

testdata:
	go run ./cmd/cli-gen/main.go testdata/cligen/generated.json internal/cligen/generated github.com/nokamoto/2pf23/internal/cligen/generated

cli:
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install ./cmd/protoc-gen-cli
	buf generate --template buf.gen.cli.yaml
	go run ./cmd/cli-gen/main.go build/cli/test.json internal/cli/generated github.com/nokamoto/2pf23/internal/cli/generated

.PHONY: all go proto cialpha mock build testdata
