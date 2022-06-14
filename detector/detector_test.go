package detector

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func TestResourceAttributes(t *testing.T) {

	// Expected resource object
	resourceAttributes := []attribute.KeyValue{
		ModuleImportPathKey.String("github.com/company/app"),
		ModulePathKey.String("/app"),
		CommitIdKey.String("123"),
		EnvironmentKey.String("dev"),
		OtherModuleImportPathKey.StringSlice(make([]string, 0)),
		OtherModulePathKey.StringSlice(make([]string, 0)),
	}
	expectedResource := resource.NewWithAttributes(semconv.SchemaURL, resourceAttributes...)

	detector := DigmaDetector{
		DeploymentEnvironment: "dev",
		CommitId:              "123",
		ModuleImportPath:      "github.com/company/app",
		ModulePath:            "/app",
	}
	resourceObj, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedResource, resourceObj, "Resource object returned is incorrect")
}

func TestShouldFailIfNoDeploymentEnvironment(t *testing.T) {
	detector := DigmaDetector{}
	_, err := detector.Detect(context.Background())
	expectedErrorMsg := "DeploymentEnvironment is required"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}

func TestShouldFailIfUnableToResolveCurrModuleInfo(t *testing.T) {
	detector := DigmaDetector{}
	_, err := detector.Detect(context.Background())
	expectedErrorMsg := "import \"\": invalid import path"
	assert.EqualErrorf(t, err, expectedErrorMsg, "Error should be: %v, got: %v", expectedErrorMsg, err)
}
