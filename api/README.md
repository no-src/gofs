# API

## Installation

Download the latest `protoc`.

`https://github.com/protocolbuffers/protobuf/releases`

Install the protocol compiler plugins for Go using the following commands:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Regenerate gRPC code.

```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto
```