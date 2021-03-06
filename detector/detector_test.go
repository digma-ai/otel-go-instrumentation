package detector

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func TestResourceAttributes(t *testing.T) {
	hostname, _ := os.Hostname()
	// Expected resource object
	resourceAttributes := []attribute.KeyValue{
		ModuleImportPathKey.String("github.com/company/app"),
		ModulePathKey.String("/app"),
		CommitIdKey.String("123"),
		EnvironmentKey.String("dev"),
		OtherModuleImportPathKey.StringSlice(make([]string, 0)),
		OtherModulePathKey.StringSlice(make([]string, 0)),
		SpanMappingPatternKey.String(""),
		SpanMappingReplacementKey.String(""),
		semconv.TelemetrySDKLanguageGo,
		semconv.HostNameKey.String(hostname),
	}
	expectedResource := resource.NewWithAttributes(semconv.SchemaURL, resourceAttributes...)

	detector := DigmaDetector{
		DigmaEnvironment: "dev",
		CommitId:         "123",
		ModuleImportPath: "github.com/company/app",
		ModulePath:       "/app",
	}
	resourceObj, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedResource, resourceObj, "Resource object returned is incorrect")
}

func TestShouldFallbackToMachineNameIfNoDeploymentEnvironment(t *testing.T) {
	hostname, _ := os.Hostname()
	resourceAttributes := []attribute.KeyValue{
		ModuleImportPathKey.String("github.com/company/app"),
		ModulePathKey.String("/app"),
		CommitIdKey.String("123"),
		EnvironmentKey.String(hostname + "[local]"),
		OtherModuleImportPathKey.StringSlice(make([]string, 0)),
		OtherModulePathKey.StringSlice(make([]string, 0)),
		SpanMappingPatternKey.String(""),
		SpanMappingReplacementKey.String(""),
		semconv.TelemetrySDKLanguageGo,
		semconv.HostNameKey.String(hostname),
	}
	expectedResource := resource.NewWithAttributes(semconv.SchemaURL, resourceAttributes...)

	detector := DigmaDetector{
		CommitId:         "123",
		ModuleImportPath: "github.com/company/app",
		ModulePath:       "/app",
	}
	resourceObj, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedResource, resourceObj, "Resource object returned is incorrect")
}

func TestShouldFailIfUnableToResolveCurrModuleInfo(t *testing.T) {
	detector := DigmaDetector{
		DigmaEnvironment: "dev",
	}
	_, err := detector.Detect(context.Background())
	expectedErrorMsg := "unable to read buildinfo. ModulePath and ModuleImportPath are required"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}
