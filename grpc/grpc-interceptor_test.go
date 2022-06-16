package grpc

import (
    "testing"
	"github.com/stretchr/testify/assert"
)

func TestMethodOnly(t *testing.T) {

	assert.Equal(t, "xyz", methodOnly("/abcd/xyz"))

}