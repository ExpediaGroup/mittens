package internal

import (
	"mittens/internal/pkg/probe"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTempFile(t *testing.T) {
	var content string = "testfile"
	var generatedFile string = CreateTempFile(content)

	// check if the file generated is exist
	var expected bool = true
	fileCheckBool, err := probe.FileExists(generatedFile)
	require.NoError(t, err)
	assert.Equal(t, expected, fileCheckBool)

	// clean up test file
	probe.DeleteFile(generatedFile)
}
