# Opentelemetry Go Instrumentation For gorilla/mux

This package provides instrumentation for additional span attributes provided on top of the [opentelmetery-instrumentation-gorilla/mux](https://go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux) package. 

In order to be able to effectively glean code-object based insights for continuous feedback and map them back in the IDE, Digma inserts additional attribute into the OTEL resource attributes. 

## Pre-requisites
* Go with Go modules enabled.
*  [opentelmetery-instrumentation-gorilla/mux](https://go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux) package

## Installing the module
```
go get -u github.com/digma-ai/otel-go-instrumentation/mux@v1.0.0
```

### Instrumenting your mux project

The Digma instrumentation depends on the otelmux opentelemetry instrumentation.

Make sure the digmamux middleware goes **after** the otelmux middleware

```go
import (
	digmamux "github.com/digma-ai/otel-go-instrumentation/mux"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)


func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(otelmux.Middleware(appName))
	router.Use(digmamux.Middleware(router))
}
```

### Additional span attributes added by this instrumentation

| Span Attribute | Example Value |
| --- | --- |
|`endpoint.function_full_name` | github.com/digma-ai/otel-sample-application-go/src/userservice/user.(*UserController).Add-fm
