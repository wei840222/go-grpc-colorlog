# go-grpc-colorlog

## usage
```go
import (
    "google.golang.org/grpc"
    grpc_colorlog "github.com/wei840222/go-grpc-colorlog"
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

myServer := grpc.NewServer(
    grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
        grpc_colorlog.StreamServerInterceptor(),
    )),
    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
        grpc_colorlog.UnaryServerInterceptor(),
    )),
)
```