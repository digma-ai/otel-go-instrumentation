package grpc

import (
	"context"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UnaryClientInterceptor returns a grpc.UnaryClientInterceptor suitable
// for use in a grpc.Dial call.
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		callOpts ...grpc.CallOption,
	) error {

		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String("endpoint.function_full_name", method))

		return nil
	}
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String("endpoint.function_full_name", info.FullMethod))

		return req, nil
	}
}
