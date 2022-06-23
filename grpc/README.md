# Opentelemetry Go Instrumentation For gRPC

This package provides instrumentation for additional span attributes provided on top of the [opentelmetery-instrumentation-grpc](https://go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc) package. 

In order to be able to effectively glean code-object based insights for continuous feedback and map them back in the IDE, Digma inserts additional attribute into the OTEL resource attributes. 

## Pre-requisites
*  Go with Go modules enabled  `version: 1.17 or above.`
*  [opentelmetery-instrumentation-grpc](https://go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc) package

## Installing the module
```
go get -u github.com/digma-ai/otel-go-instrumentation/grpc@v1.0.1
```

### Instrumenting your gRPC project

The Digma instrumentation depends on the otelgrpc opentelemetry instrumentation.

Make sure the digmagrpc middleware goes **after** the otelgrpc middleware

```go
import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	digmagrpc "github.com/digma-ai/otel-go-instrumentation/grpc"
)


func main() {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			otelgrpc.UnaryServerInterceptor(),
			digmagrpc.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			otelgrpc.StreamServerInterceptor(),
			digmagrpc.StreamServerInterceptor(),
		)),
	)
}
```

### Additional span attributes added by this instrumentation

| Span Attribute | Example Value |
| --- | --- |
|`endpoint.function_full_name` | main.(*server).SayHello
