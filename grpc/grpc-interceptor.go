package grpc

import (
	"context"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	//
	//"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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

// StreamClientInterceptor returns a grpc.StreamClientInterceptor suitable
// for use in a grpc.Dial call.
func StreamClientInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		callOpts ...grpc.CallOption,
	) (grpc.ClientStream, error) {

		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String("endpoint.function_full_name", method))

		// s, err := streamer(ctx, desc, cc, method, callOpts...)
		// if err != nil {
		// 	return s, err
		// }
		// stream := otelgrpc.wrapClientStream(ctx, s, desc)

		//TODO: check if return nil as grpc.ClientStream is valid...
		return nil, nil
	}
}

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor suitable
// for use in a grpc.NewServer call.
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

// StreamServerInterceptor returns a grpc.StreamServerInterceptor suitable
// for use in a grpc.NewServer call.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {

		ctx := ss.Context()
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String("endpoint.function_full_name", info.FullMethod))

		return nil
	}
}