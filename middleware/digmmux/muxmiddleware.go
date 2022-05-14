package digmmux

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

type traceware struct {
}

func Middleware(service string, opts ...otelmux.Option) mux.MiddlewareFunc {
	return func(handler http.Handler) http.Handler {
		return traceware{}
	}
}

func (tw traceware) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
