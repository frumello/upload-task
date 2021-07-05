# Upload Image Task

## Tech/Framework used

<b>Scaffolding</b>
- [GO MOD INIT] (https://blog.golang.org/using-go-modules)
- [Protocol Buffers] (https://developers.google.com/protocol-buffers)

<b>Built with</b>
- [Go 1.6] (http://www.oracle.com/technetwork/java/javase/overview/index.html)
- [gRPC-Go] (https://pkg.go.dev/google.golang.org/grpc)
- [Go support for Protocol Buffers] (https://pkg.go.dev/google.golang.org/protobuf)
- [Testify - Thou Shalt Write Tests] (https://github.com/stretchr/testify)

## Packaging and running

Commands below must be executed from the repository root

### Build and Run upload-server

`make run` or `docker compose up --build -d`

### Install grpc-upload (client)

`make install-client`

*you must have GOPATH in your PATH env

### Stop upload-server when you finish

`make stop` or `docker-compose down`

## To upload a file
Upload should be possible using the syntax:

`grpc-upload ./cute_kitten.jpg`

Download the file should be possible in the url:

`http://localhost:8888/cute_kitten.jpg`