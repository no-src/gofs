# Internal gRPC API

## Install Dependent Tools

Download the latest `protoc`.

`https://github.com/protocolbuffers/protobuf/releases`

Install the protocol compiler plugins for Go using the following commands.

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Regenerate gRPC code

Firstly, you should switch to the root path of the repository like the following command.

```bash
cd $GOPATH/src/github.com/no-src/gofs
```

Then you can use the following command to regenerate gRPC code.

```bash
protoc --go_out=../../../ --go_opt=paths=import --go-grpc_out=../../../ --go-grpc_opt=paths=import api/proto/*.proto
```