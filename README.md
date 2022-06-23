# Opentelemetry Go Instrumentation Digma
[![Tests](https://github.com/digma-ai/otel-go-instrumentation/actions/workflows/ci.yaml/badge.svg)](https://github.com/digma-ai/otel-go-instrumentation/actions?query=workflow%3Abuild_and_test+branch%3Amain)
[![Docs](https://godoc.org/go.opentelemetry.io/contrib?status.svg)][goref-url]

This package provides instrumentation to make it easy to set up Digma to work along with your OpenTelemetry instrumentation.

In order to be able to effectively glean code-object based insights for continuous feedback and map them back in the IDE, Digma inserts additional attribute into the OTEL resource attributes. 

## Pre-requisites
*  Go with Go modules enabled.


## Installing the module
```
go get -u github.com/digma-ai/otel-go-instrumentation@v1.0.8
```


## Usage

### Set up
- [Initiallizing opentelemetry resource](#initiallizing-opentelemetry-resource)
- [Framework instrumentation](#framework-instrumentation)


### Initiallizing opentelemetry resource

If you have an existing OpenTelemtry instrumentaiton set up, simply use the DigmaDetector object to create a `Resource `object and merge it with your resource to import all of the needed attributes. 

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
> Ensure required service name (semconv.ServiceNameKey) is set on resource.

Now your TracerProvider will have the following resource attributes and attach them to new spans:

| Resource Attribute | Example Value |
| --- | --- |
|`deployment.environment` | Production [or hostname if not set]
|`scm.commit.id` | 07e239f2f3d8adc12566eaf66e0ad670f36202b5 [OPTIONAL]
|`code.module.importpath` | github.com/digma-ai/otel-go-instrumentation
|`code.module.path` | /build/work/otel-go-instrumentation



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



### The Digma instrumentation options

| Options | Input Type  | Attribute Key | Description | Default |
| --- | --- | --- | --- | --- |
| `DeploymentEnvironment` | `string` | deployment.environment |  The Environment describes where the running process is deployed. (e.g production, staging, ci) | `os.Hostname()`
| `CommitId` | `string`  | scm.commit.id | The specific commit identifier of the running code. | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()`|
`ModuleImportPath` | `string` | code.module.importpath | Module canonical name | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()`|
`ModulePath` | `string` | code.module.path | workspace(application) physical path | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()` |
`OtherModulesImportPath` | `[] string` | code.othermodule.importpath | Specify additional satellite or infra modules to track | None |
` **Internal** ` | `[] string` | code.othermodule.path | physical paths of  `OtherModulesImportPath` option | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()`|

### Framework instrumentation

The following instrumentation packages provides instrumentation for additional span attributes provided on top of the opentelmetery-instrumentation packages.

In order to be able to effectively glean code-object based insights for continuous feedback and map them back in the IDE, Digma inserts additional attribute into the OTEL resource attributes.

| Instrumentation Package |
| :---------------------: |
| [github.com/digma-ai/otel-go-instrumentation/echo](./echo) |
| [github.com/digma-ai/otel-go-instrumentation/grpc](./grpc)|  
| [github.com/digma-ai/otel-go-instrumentation/mux](./mux) | 


The [OpenTelemetry registry](https://opentelemetry.io/registry/) is the best place to discover instrumentation packages.
It will include packages outside of this project.


[goref-url]: https://pkg.go.dev/github.com/digma-ai/otel-go-instrumentation
