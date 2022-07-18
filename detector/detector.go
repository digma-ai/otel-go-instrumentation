package detector

import (
	"context"
	"errors"
	"go/build"
	"os"
	"runtime/debug"
	"strings"

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
	ModuleImportPath       string //module canonical name
	ModulePath             string // workspace(application) physical path
	SpanMappingPattern     string
	SpanMappingReplacement string
}

// compile time assertion that Digma implements the resource.Detector interface.
var _ resource.Detector = (*DigmaDetector)(nil)

func (d *DigmaDetector) Detect(ctx context.Context) (*resource.Resource, error) {
	deploymentEnvironment := strings.TrimSpace(d.DeploymentEnvironment)
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	if deploymentEnvironment == "" {
		deploymentEnvironment = hostname + "[local]"
	}

	attributes := []attribute.KeyValue{
		EnvironmentKey.String(deploymentEnvironment)}

	moduleImportPath := strings.TrimSpace(d.ModuleImportPath)
	modulePath := strings.TrimSpace(d.ModulePath)
	commitId := strings.TrimSpace(d.CommitId)

	//module name and path explicit defined by user
	if moduleImportPath != "" && modulePath != "" {
		attributes = append(attributes, ModuleImportPathKey.String(moduleImportPath))
		attributes = append(attributes, ModulePathKey.String(modulePath))
	} else if moduleImportPath != "" {
		return nil, errors.New("ModulePath is required")
	} else if modulePath != "" {
		return nil, errors.New("ModuleImportPath is required")
	}

	if bi, ok := debug.ReadBuildInfo(); ok {
		if commitId != "" {
			attributes = append(attributes, CommitIdKey.String(d.CommitId))
		}
		// else {
		// 	for _, setting := range bi.Settings {
		// 		if setting.Key == "vcs.revision" {
		// 			attributes = append(attributes, CommitIdKey.String(setting.Value))
		// 			break
		// 		}
		// 	}
		// }
		if moduleImportPath == "" && modulePath == "" {
			attributes = append(attributes, ModuleImportPathKey.String(bi.Main.Path)) //module path
			imported, err := build.Default.Import(bi.Path, ".", build.FindOnly)
			if err != nil {
				return nil, err
			} else {
				attributes = append(attributes, ModulePathKey.String(imported.Root))
			}
		}
	}

	var otherModulesImportPath []string
	var otherModulesPath []string

	for i := 0; i < len(d.OtherModulesImportPath); i++ {

		imported, err := build.Default.Import(d.OtherModulesImportPath[i], modulePath, build.FindOnly)
		if err != nil {
			return nil, err
		}
		otherModulesImportPath = append(otherModulesImportPath, imported.ImportPath)
		otherModulesPath = append(otherModulesPath, imported.Root)
	}
	attributes = append(attributes, OtherModuleImportPathKey.StringSlice(otherModulesImportPath))
	attributes = append(attributes, OtherModulePathKey.StringSlice(otherModulesPath))

	attributes = append(attributes, SpanMappingPatternKey.String(d.SpanMappingPattern))
	attributes = append(attributes, SpanMappingReplacementKey.String(d.SpanMappingReplacement))
	attributes = append(attributes, semconv.HostNameKey.String(hostname))

	attributes = append(attributes, semconv.TelemetrySDKLanguageGo)
	return resource.NewWithAttributes(semconv.SchemaURL, attributes...), nil
}
