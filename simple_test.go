package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//TestSomething Test method 1
func TestSomething(t *testing.T) {
	ass := assert.New(t)
	// assert equality
	ass.Equal(123, 123, "they should be equal")
}

//TestAverage Test method 2
func TestAverage(t *testing.T) {
	ass := assert.New(t)
	// assert inequality
	ass.NotEqual(123, 456, "they should not be equal")
}
