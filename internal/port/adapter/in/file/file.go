// Package file provides functionality
// for getting input data from a file located in a file system.
package file

import (
	"io"
	"os"
)

type Input struct {
	file *os.File
}

func NewInput(filePath string) (*Input, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &Input{
		file,
	}, nil
}

// GetData function returns an io.Reader that can be used to read the contents of the file
func (i *Input) GetData() io.Reader {
	return i.file
}

func (i *Input) Stop() error {
	return i.file.Close()
}
