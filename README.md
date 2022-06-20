# opentelemetry-go-instrumentation-digma
[![Tests](https://github.com/digma-ai/otel-go-instrumentation/actions/workflows/ci.yaml/badge.svg)](https://github.com/digma-ai/otel-go-instrumentation/actions/workflows/ci.yaml)

This package provides instrumentation helpers and tools to make it easy to set up Digma to work along with your OpenTelemetry instrumentation.

In order to be able to effectively glean code-object based insights for continuous feedback and map them back in the IDE, Digma inserts additional attribute into the OTEL resource attributes. 


## Installing the module
```
go get -u github.com/digma-ai/otel-go-instrumentation@v1.0.1
```


## Usage

If you have an existing OpenTelemtry instrumentaiton set up, simply use the DigmaDetector object to create a `Resource `object and merge it with your resource to import all of the needed attributes. 

```go
import (
	"github.com/digma-ai/otel-go-instrumentation/detector"
)

resource.WithDetectors(
    &detector.DigmaDetector{
        DeploymentEnvironment:       "Production",
        CommitId:                    "", //optional
    },
)
```
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



## The Digma instrumentation options

| Options | Input Type  | Attribute Key | Description | Default |
| --- | --- | --- | --- | --- |
| `DeploymentEnvironment` | `string` | deployment.environment |  The Environment describes where the running process is deployed. (e.g production, staging, ci) | `os.Hostname()`
| `CommitId` | `string`  | scm.commit.id | The specific commit identifier of the running code. | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()`|
`ModuleImportPath` | `string` | code.module.importpath | Module canonical name | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()`|
`ModulePath` | `string` | code.module.path | workspace(application) physical path | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()` |
`OtherModulesImportPath` | `[] string` | code.othermodule.importpath | Specify additional satellite or infra modules to track | None |
` **Internal** ` | `[] string` | code.othermodule.path | physical paths of  `OtherModulesImportPath` option | The instrumentation will attempt to read this variable from `debug.ReadBuildInfo()`


[goref-url]: https://pkg.go.dev/github.com/digma-ai/otel-go-instrumentation