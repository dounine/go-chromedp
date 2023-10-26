package test

import (
	"fmt"
	assert2 "github.com/stretchr/testify/assert"
	"testing"
)

func TestFormat(t *testing.T) {
	assert := assert2.New(t)
	assert.Equal(fmt.Sprintf("%s/%d", "1", 1), "1/1")
}
