package grpc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethodOnly(t *testing.T) {

	assert.Equal(t, "Xyz", methodOnly("/Abcd/Xyz"))
}

func TestBuildMethodFqnGoodPath(t *testing.T) {

	impl := &impl4tests{}
	methodFqn, err := buildMethodFqn(impl, "/api/DoSomething")
	assert.Nil(t, err, "err should be nil")
	assertEndsWith(t, methodFqn, "(*impl4tests).DoSomething")
}

func assertEndsWith(t *testing.T, entireValue string, expectedSuffix string) {
	assert.True(t, strings.HasSuffix(entireValue, expectedSuffix), "'"+entireValue+"' has no expected suffix '"+expectedSuffix+"'")
}

type iface4tests interface {
	DoSomething()
}

type impl4tests struct {
	iface4tests
}

func (impl *impl4tests) DoSomething() {
	// nothing, just signature
}
