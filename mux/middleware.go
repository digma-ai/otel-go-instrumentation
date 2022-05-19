package mux

import (
	"errors"
	"net/http"
	"reflect"
	"runtime"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Middleware(router *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			name, err := getHandlerName(router, r)
			if err == nil {
				span := trace.SpanFromContext(r.Context())
				span.SetAttributes(attribute.String("endpoint.function_full_name", name))
			}
			next.ServeHTTP(w, r)
		})
	}
}

func getHandlerName(router *mux.Router, request *http.Request) (string, error) {
	var match mux.RouteMatch
	var handler http.Handler
	if router.Match(request, &match) {
		handler = match.Route.GetHandler()
		name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		return name, nil
	}
	return "", errors.New("Handler not found")
}
