# Opentelemetry Go Instrumentation For Echo

This package provides instrumentation for additional span attributes provided on top of the [opentelmetery-instrumentation-echo](https://go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho) package. 

In order to be able to effectively glean code-object based insights for continuous feedback and map them back in the IDE, Digma inserts additional attribute into the OTEL resource attributes. 

## Pre-requisites
*  Go with Go modules enabled.
*  [opentelmetery-instrumentation-echo](https://go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho) package

## Installing the module
```
go get -u github.com/digma-ai/otel-go-instrumentation/echo@v1.0.1
```

### Instrumenting your echo project

The Digma instrumentation depends on the echo opentelemetry instrumentation.

Make sure the digmaecho middleware goes **after** the otelecho middleware

```go
import (
	digmaecho "github.com/digma-ai/otel-go-instrumentation/echo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)


func main() {
	r := echo.New()
	r.Use(otelecho.Middleware(appName))
	r.Use(digmaecho.Middleware())
	r.GET("/", index)
}
```

### Additional span attributes added by this instrumentation

| Span Attribute | Example Value |
| --- | --- |
|`endpoint.function_full_name` | github.com/digma-ai/otel-sample-application-go/src/authservice/auth.(*AuthController).Authenticate-fm
