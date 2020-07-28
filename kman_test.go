package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilteringFiles(t *testing.T) {
	files := []string{"green apple", "red apple", "yellow apple", "yellow banana", "brown banana"}

	t.Run("simple filtering", func(t *testing.T) {
		filteredFiles, _ := filterFiles(files, "apple")

		assert.Equal(t, []string{"green apple", "red apple", "yellow apple"}, filteredFiles)
	})

	t.Run("filtering with wildcard", func(t *testing.T) {
		filteredFiles, _ := filterFiles(files, "*ana")

		assert.Equal(t, []string{"yellow banana", "brown banana"}, filteredFiles)
	})

	t.Run("filtering with space", func(t *testing.T) {
		filteredFiles, _ := filterFiles(files, " apple")

		assert.Equal(t, []string{"green apple", "red apple", "yellow apple"}, filteredFiles)
	})

	t.Run("filtering with single character wildcard", func(t *testing.T) {
		filteredFiles, _ := filterFiles(files, "a??le")

		assert.Equal(t, []string{"green apple", "red apple", "yellow apple"}, filteredFiles)
	})
}
