package detector

import (
	"context"
	"fmt"
	"go/build"
	"runtime/debug"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

const (
	CommitIdKey         = attribute.Key("scm.commit.id")
	ModuleImportPathKey = attribute.Key("code.module.importpath")
	ModulePathKey       = attribute.Key("code.module.path")

	OtherModuleImportPathKey = attribute.Key("code.othermodule.importpath")
	OtherModulePathKey       = attribute.Key("code.othermodule.path")

	EnvironmentKey = semconv.DeploymentEnvironmentKey

	SpanMappingPatternKey     = attribute.Key("digma.span_mapping_pattern")
	SpanMappingReplacementKey = attribute.Key("digma.span_mapping_replacement")
)

type DigmaDetector struct {
	DeploymentEnvironment  string
	CommitId               string
	OtherModulesImportPath []string
	SpanMappingPattern     string
	SpanMappingReplacement string
}

// compile time assertion that Digma implements the resource.Detector interface.
var _ resource.Detector = (*DigmaDetector)(nil)

func (d *DigmaDetector) Detect(ctx context.Context) (*resource.Resource, error) {

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

		attributes = append(attributes, ModuleImportPathKey.String(bi.Main.Path))
		imported, err := build.Default.Import(bi.Main.Path, ".", build.FindOnly)
		if err != nil {
			panic(err)
		} else {
			attributes = append(attributes, ModulePathKey.String(imported.Root))
		}

		var otherModulesImportPath []string
		var otherModulesPath []string

		for i := 0; i < len(d.OtherModulesImportPath); i++ {

			imported, err := build.Default.Import(d.OtherModulesImportPath[i], ".", build.FindOnly)
			if err != nil {
				panic(err)
			}
			otherModulesImportPath = append(otherModulesImportPath, imported.ImportPath)
			otherModulesPath = append(otherModulesPath, imported.Root)
		}
		attributes = append(attributes, OtherModuleImportPathKey.StringSlice(otherModulesImportPath))
		attributes = append(attributes, OtherModulePathKey.StringSlice(otherModulesPath))

		attributes = append(attributes, SpanMappingPatternKey.String(d.SpanMappingPattern))
		attributes = append(attributes, SpanMappingReplacementKey.String(d.SpanMappingReplacement))
	}

	var err error

	fmt.Println("digma attributes:")
	for _, attr := range attributes {
		fmt.Printf("%s=%s\n", attr.Key, attr.Value.Emit())
	}

	return resource.NewSchemaless(attributes...), err
}
