API_GO_FILES = $(shell find pkg/api -name '*grpc.pb.go')
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

define versioning
	echo $1
	$1 --version > build/tools/$1
	sha1sum build/tools/* | sha1sum > build/.tools
endef

fast:
	go test ./...

all: proto mock testdata gen go

$(GOBIN)/gofumpt:
	go install mvdan.cc/gofumpt@latest
	$(call versioning,gofumpt)

$(GOBIN)/staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	$(call versioning,staticcheck)

$(GOBIN)/protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(call versioning,protoc-gen-go)

$(GOBIN)/protoc-gen-go-grpc:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(call versioning,protoc-gen-go-grpc)

$(GOBIN)/buf:
	go install github.com/bufbuild/buf/cmd/buf@latest
	$(call versioning,buf)

$(GOBIN)/mockgen:
	go install github.com/golang/mock/mockgen@latest
	$(call versioning,mockgen)

$(GOBIN)/ko:
	go install github.com/google/ko@latest
	$(call versioning,ko)

go: $(GOBIN)/gofumpt $(GOBIN)/staticcheck $(GOBIN)/mockgen
	go generate ./...
	gofumpt -l -w .
	go vet ./...
	staticcheck ./...
	go test ./... -race -covermode=atomic -coverprofile=coverage.out
	go mod tidy

proto: $(GOBIN)/protoc-gen-go $(GOBIN)/protoc-gen-go-grpc $(GOBIN)/buf
	buf lint --config build/buf/buf.yaml --error-format=json
	buf format --config build/buf/buf.yaml -w
	buf generate --template build/buf/buf.gen.yaml

cialpha:
	go run ./build/ci/test/main.go

mock: $(GOBIN)/mockgen
	$(foreach file,$(API_GO_FILES),mockgen -source $(file) -destination internal/mock/$(file))

.PHONY: build
build: $(GOBIN)/ko
	$(foreach command,$(COMMANDS),$(call ko,$(command)))

.PHONY: testdata
testdata:
	go run ./cmd/cli-gen/main.go testdata/cligen/generated.json internal/cligen/generated github.com/nokamoto/2pf23/internal/cligen/generated
	go run ./cmd/server-gen/main.go testdata/servergen internal/servergen/generated --mock

gen: $(GOBIN)/buf
	go install ./cmd/protoc-gen-cli
	go install ./cmd/protoc-gen-server
	buf generate --template build/buf/buf.gen.local.yaml
	go run ./cmd/cli-gen/main.go build/cli/test.json internal/cli/generated github.com/nokamoto/2pf23/internal/cli/generated
	go run ./cmd/server-gen/main.go build/server internal/server/generated

tilt:
	curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
