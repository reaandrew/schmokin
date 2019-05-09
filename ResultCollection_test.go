package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultCollection(t *testing.T) {

	t.Run("Success returns true", func(t *testing.T) {
		var collection = ResultCollection{}
		collection = append(collection, Result{
			Success: true,
		})
		assert.True(t, collection.Success())
	})

	t.Run("Success returns false", func(t *testing.T) {
		var collection = ResultCollection{}
		collection = append(collection, Result{
			Success: false,
		})
		assert.False(t, collection.Success())
	})
}
