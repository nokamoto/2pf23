# Design: Protocol Buffers Code Generation

## Overview
```mermaid
flowchart LR
    api[api proto files] --buf--> buildcli[build/cli json files]
    buildcli --> gocli[internal/cli/generated go files]
    api --buf--> buildserver[build/server json files]
    buildserver --> goserver[internal/server/generated go files]
```

The `api` directory contains the protocol buffer files for the API. For each [standard methods](https://cloud.google.com/apis/design/standard_methods) defined in the API, corresponding client and server Go codes are generated.

For example, the `Create` method is defined in [ke.proto](../../api/ke/v1alpha/ke.proto):

```proto
service KeService {
  rpc CreateCluster(CreateClusterRequest) returns (Cluster);
}
```

For the `Create` method, the following [cobra](https://github.com/spf13/cobra) implementation is generated in [createcluster.go](../../internal/cli/generated/ke/v1alpha/cluster/createcluster.go):

```go
cmd := &cobra.Command{
    Use: "create",
    RunE: func(cmd *cobra.Command, args []string) error {
        // ...
    },
}
```

And the following gRPC service implementation is generated in [service.go](../../internal/server/generated/ke/v1alpha/service.go):

```go
type service struct {
	v1alpha.UnimplementedKeServiceServer
}

func (s *service) CreateCluster(ctx context.Context, req *v1alpha.CreateClusterRequest) (*v1alpha.Cluster, error) {
    // ...
}
```

## Code generation for CLI
| input | command | output |
| --- | --- | --- |
| [api](../../api) | [protoc-gen-cli](../../cmd/protoc-gen-cli/) | [build/cli](../../build/cli) |
| [build/cli](../../build/cli) | [cli-gen](../../cmd/cli-gen/) | [internal/cli/generated](../../internal/cli/generated) |

```bash
buf generate --template buf.gen.local.yaml
```

The `protoc-gen-cli` plugin is used to generate JSON file(s) for each CLI command. The JSON file contains the command name, description, flags, etc.. The `inhouse.v1.Package` proto message is used to define the CLI commands.

See [generated.json](../../testdata/cligen/generated.json) for an example of the generated JSON file.

```bash
go run ./cmd/cli-gen/main.go build/cli/test.json internal/cli/generated github.com/nokamoto/2pf23/internal/cli/generated
```

The `cli-gen` command is used to generate the Go code for the CLI commands. The generated code is placed in the `internal/cli/generated` directory.

## Code generation for server
| input | command | output |
| --- | --- | --- |
| [api](../../api) | [protoc-gen-server](../../cmd/protoc-gen-server/) | [build/server](../../build/server) |
| [build/server](../../build/server) | [server-gen](../../cmd/server-gen/) | [internal/server/generated](../../internal/server/generated) |

```bash
buf generate --template buf.gen.local.yaml
```

The `protoc-gen-server` plugin is used to generate JSON file(s) for each gRPC service. The JSON file contains the service name, methods, input/output messages, etc.. The `inhouse.v1.Service` proto message is used to define the gRPC services.

See [ke.json](../../testdata/servergen/ke.json) for an example of the generated JSON file.

```bash

```bash
go run ./cmd/server-gen/main.go build/server internal/server/generated
```

The `server-gen` command is used to generate the Go code for the gRPC services. The generated code is placed in the `internal/server/generated` directory.
