package probe

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// I think we can use this test to be FileExist unit test function as well.
func TestWriteFile(t *testing.T) {
	var fileName string = "testfile"
	var expected bool = true
	WriteFile(fileName)

	// check is file exist
	result, err := FileExists(fileName)
	require.NoError(t, err)

	// delete file afterwards
	DeleteFile(fileName)

	assert.Equal(t, expected, result)
}

func TestDeleteFile(t *testing.T) {
	var fileName string = "testfile"
	var expected bool = false
	WriteFile(fileName)

	// delete file afterwards
	DeleteFile(fileName)

	// check is file exist
	result, err := FileExists(fileName)
	require.NoError(t, err)

	assert.Equal(t, expected, result)
}
