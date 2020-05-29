package gogen

import (
	"github.com/pkg/errors"
	"path/filepath"
)

type File interface {
	Path() string
	Data
}

type Data interface {
	Bytes() []byte
}

type FileWriter interface {
	Write(file File) error
}

type Files []File

// should never be called. Used to implement File interface
func (d Files) Bytes() []byte {
	return nil
}

// should never be called. Used to implement File interface
func (d Files) String() string {
	return ""
}

// should never be called. Used to implement File interface
func (d Files) Path() string {
	return ""
}

func Write(file File) error {
	if w, ok := file.(FileWriter); ok {
		return w.Write(file)
	}

	if err := makeDir(filepath.Dir(file.Path())); err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	if err := createFile(file.Path(), file.Bytes()); err != nil {
		return errors.Wrapf(err, "failed to write %s", file.Path())
	}
	return nil
}
