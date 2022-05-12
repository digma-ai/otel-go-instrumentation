package detector

import (
	"context"
	"fmt"
	"runtime/debug"

	//"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

const (
	CommitIdKey    = attribute.Key("scm.commit.id")
	ModuleKey      = attribute.Key("code.module.path")
	EnvironmentKey = semconv.DeploymentEnvironmentKey
)

type Digma struct {
	DeploymentEnvironment string
	CommitId              string
}

// compile time assertion that Digma implements the resource.Detector interface.
var _ resource.Detector = (*Digma)(nil)

func (d *Digma) Detect(ctx context.Context) (*resource.Resource, error) {

	attributes := []attribute.KeyValue{
		EnvironmentKey.String(d.DeploymentEnvironment)}

	if bi, ok := debug.ReadBuildInfo(); ok {
		if len(d.CommitId) > 0 {
			attributes = append(attributes, CommitIdKey.String(d.CommitId))
		} else {
			for _, setting := range bi.Settings {
				if setting.Key == "vcs.revision" {
					attributes = append(attributes, CommitIdKey.String(setting.Value))
					break
				}
			}
		}

		//The main package path
		attributes = append(attributes, ModuleKey.String(bi.Main.Path))
	}

	var err error

	fmt.Println("digma attributes:")
	for _, attr := range attributes {
		fmt.Printf("%s=%s\n", attr.Key, attr.Value.AsString())
	}

	return resource.NewSchemaless(attributes...), err
}
