package files

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFileWithContent(t *testing.T) {
	fileName := "testfile.txt"
	content := "This is a test."

	defer os.Remove(fileName)

	err := CreateFileWithContent(fileName, content)

	assert.NoError(t, err, "expected no error when creating file with content")

	fileData, err := os.ReadFile(fileName)
	assert.NoError(t, err, "expected no error when reading the file")
	assert.Equal(t, content, string(fileData), "file content should match the input content")
}
