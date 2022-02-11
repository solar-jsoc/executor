package executor

import (
	"os"
)

type file struct {
	path string
	file *os.File
}

type inputFile file

func (i *inputFile) Read(b []byte) (n int, err error) {
	if i.file == nil {
		f, err := os.Open(i.path)
		if err != nil {
			return 0, err
		}
		i.file = f
	}

	return i.file.Read(b)
}

func (i *inputFile) Close() error {
	if i.file != nil {
		return i.file.Close()
	}

	return nil
}

type outputFile file

func (o *outputFile) Write(b []byte) (n int, err error) {
	if o.file == nil {
		f, err := os.Create(o.path)
		if err != nil {
			return 0, err
		}
		o.file = f
	}
	return o.file.Write(b)
}

func (o *outputFile) Close() error {
	if o.file != nil {
		return o.file.Close()
	}

	return nil
}
