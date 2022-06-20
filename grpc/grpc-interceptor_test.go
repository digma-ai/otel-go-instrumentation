package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	stack_assert "github.com/subchen/go-stack/assert"
)

func TestMethodOnly(t *testing.T) {

	assert.Equal(t, "Xyz", methodOnly("/Abcd/Xyz"))
}

func TestBuildMethodFqnGoodPath(t *testing.T) {

	impl := &impl4tests{}
	methodFqn, err := buildMethodFqn(impl, "/api/DoSomething")
	if assert.NoError(t, err) {
		stack_assert.HasSuffix(t, methodFqn, "(*impl4tests).DoSomething")
	}
}

func TestBuildMethodFqnBadPath(t *testing.T) {

	impl := &impl4tests{}
	_, err := buildMethodFqn(impl, "/api/Do123")
	assert.Error(t, err)
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
