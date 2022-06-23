# Opentelemetry Go Instrumentation Digma
[![Tests](https://github.com/digma-ai/otel-go-instrumentation/actions/workflows/ci.yaml/badge.svg)](https://github.com/digma-ai/otel-go-instrumentation/actions?query=workflow%3Abuild_and_test+branch%3Amain)
[![Docs](https://godoc.org/go.opentelemetry.io/contrib?status.svg)][goref-url]

This package provides instrumentation to make it easy to set up Digma to work along with your OpenTelemetry instrumentation.

In order to be able to effectively glean code-object based insights for continuous feedback and map them back in the IDE, Digma inserts additional attribute into the OTEL resource attributes. 

## Pre-requisites
*  Go with Go modules enabled  `version: 1.17 or above.`

## Installing the module
```
go get -u github.com/digma-ai/otel-go-instrumentation@v1.0.8
```


## Usage

### Set up
- [Instrumenting your OpenTelemetry resource](#instrumenting-your-opentelemetry-resource)
- [Adding instrumentation for specific server frameworks](#adding-instrumentation-for-specific-server-frameworks)
- [Exporting trace data to Digma](#exporting-trace-data-to-digma)
- [Framework instrumentation](#framework-instrumentation)


### Instrumenting your OpenTelemetry resource

If you have an existing OpenTelemtry instrumentaiton set up, simply use the DigmaDetector object to create a `Resource ` and merge it with your existing OTEL resource to import all of the needed attributes. 

```go
import (
	"github.com/digma-ai/otel-go-instrumentation/detector"
)

res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends and mandatory for digma backend
			semconv.ServiceNameKey.String(serviceName),
		),

		resource.WithDetectors(
			&detector.DigmaDetector{
				DeploymentEnvironment: "Production",
				CommitId:              "", //optional
			},
		))
```

This will add a few additional attributes to the OTEL `Resource` component needed by Digma to map the information back to the code.

### Adding instrumentation for specific server frameworks

Digma can also generate specifc insights based on the service framework you're using. To do that, we can add  a simple middleware that will save the contexual information needed to map the tracing data to the underlying code in the Span attributes.

Follow the steps in the below links to add Digma's middlware, based on your server framework:

* [github.com/digma-ai/otel-go-instrumentation/echo](./echo)
* [github.com/digma-ai/otel-go-instrumentation/grpc](./grpc)
* [github.com/digma-ai/otel-go-instrumentation/mux](./mux)  

For example, here  is how you would use Digma's middlware with [Echo](https://github.com/labstack/echo) along with the standard OTEL middlware:

```go
func main() {
	r := echo.New()
	r.Use(otelecho.Middleware(appName))
	r.Use(digmaecho.Middleware())
	r.GET("/", index)
}
```

### Exporting trace data to Digma

First, you need to have a Digma backend up and running. You can follow the instructions in the [Digma project repository](https://github.com/digma-ai/digma) to quickly get set up using Docker.

You can use a standard OTLP exporter to the Digma collector for local deployments:

```go
import (
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"

)

traceClient := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint("localhost:5050"),
	)
```

Alternative, if you're already using a collector component you can simply modify its configuration file:

```yaml
exporters:
...
otlp/digma:
    endpoint: "localhost:5050"
    tls:
      insecure: true
service:
  pipelines:
    traces:
      exporters: [otlp/digma, ...]
```

In both cases, set the endpoint value to the URL of the Digma backend.

That's it! You should be good to go.

### The Digma instrumentation options

The following options can be set as a part of the OTEL Resource instrumentation. All of these are optional, but can help provide more context to the colleced traces:

| Options | Input Type  | Attribute Key | Description | Default |
| --- | --- | --- | --- | --- |
| `DeploymentEnvironment` | `string` | deployment.environment |  The Environment describes where the running process is deployed. (e.g production, staging, ci) | If no depployment environment is provided, we'll assume this is a local deployment env and mark it using the local hostname. It will be visible to that machine only.
| `CommitId` | `string`  | scm.commit.id | The specific commit identifier of the running code. | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()`|
`ModuleImportPath` | `string` | code.module.importpath | Module canonical name | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()`|
`ModulePath` | `string` | code.module.path | workspace(application) physical path | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()` |
`OtherModulesImportPath` | `[] string` | code.othermodule.importpath | Specify additional satellite or infra modules to track | None |
` **Internal** ` | `[] string` | code.othermodule.path | physical paths of  `OtherModulesImportPath` option | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()`|

### Additional instrumentation

The more instrumentation you add to your project, the more insights Digma will be able to provide.

The [OpenTelemetry registry](https://opentelemetry.io/registry/) is the best place to discover instrumentation packages.


[goref-url]: https://pkg.go.dev/github.com/digma-ai/otel-go-instrumentation


