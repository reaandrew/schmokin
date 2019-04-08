package main_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func CheckError(t *testing.T, err error) {
	assert.Nil(t, err)
}
