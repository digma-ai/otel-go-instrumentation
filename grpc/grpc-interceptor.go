package grpc

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"google.golang.org/grpc"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

//
// Digma's grpc-interceptor:
// should be chained after the original otelgrpc interceptors.
//
// usage for example :
//  import (
//     	"google.golang.org/grpc"
//      "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
//  	digma_grpc "github.com/digma-ai/otel-go-instrumentation/grpc"
//   	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
//  )
//
//	server := grpc.NewServer(
//	    grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
//	       otelgrpc.UnaryServerInterceptor(),
//	       digmagrpc.UnaryServerInterceptor(),
//      )),
//      grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
//	      otelgrpc.StreamServerInterceptor(),
//	      digma_grpc.StreamServerInterceptor(),
//      )),
//  )
//


// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor suitable
// for use in a grpc.NewServer call.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		methodFqn := buildMethodFqn(info.Server, info.FullMethod)
		//TODO: debug line, remove it
		fmt.Println("methodFqn:", methodFqn)

		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String("endpoint.function_full_name", methodFqn))

		// standard call of interceptor
		resp, err := handler(ctx, req)
		return resp, err
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

		methodFqn := buildMethodFqn(srv, info.FullMethod)
		//TODO: debug line, remove it
		fmt.Println("methodFqn:", methodFqn)

		ctx := ss.Context()
		span := trace.SpanFromContext(ctx)
		span.SetAttributes(attribute.String("endpoint.function_full_name", methodFqn))

		// standard call of interceptor
		err := handler(srv, ss)
		return err
	}
}

func methodOnly(fullMethod string) string {
	ix := strings.LastIndex(fullMethod, "/")

	return fullMethod[ix+1:]
}

func fqnOfServiceWithMethod(srv interface{}, methodName string) string {
	typeOfService := reflect.TypeOf(srv)
	methodRef, ok := typeOfService.MethodByName(methodName)
	if !ok {
		panic("cannot find method name '" + methodName + "' for service '" + typeOfService.Name() + "'")
	}
	methodFunc4pc := runtime.FuncForPC(methodRef.Func.Pointer())
	fqn := methodFunc4pc.Name()
	return fqn
}

func buildMethodFqn(srv interface{}, fullMethod string) string {
	methodName := methodOnly(fullMethod)
	fqn := fqnOfServiceWithMethod(srv, methodName)

	return fqn
}
