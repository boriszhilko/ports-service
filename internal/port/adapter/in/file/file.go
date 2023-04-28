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

func (i *Input) GetData() io.Reader {
	return i.file
}

func (i *Input) Stop() error {
	return i.file.Close()
}
