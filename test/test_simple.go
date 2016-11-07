package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	ass := assert.New(t)

	// assert equality
	ass.Equal(123, 123, "they should be equal")

	// assert inequality
	ass.NotEqual(123, 456, "they should not be equal")
}

func TestAverage(t *testing.T) {

}
