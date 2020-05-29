package gogen

import (
	"fmt"
	"github.com/pkg/errors"
	"path/filepath"
)

type List interface {
	SetTemplateDataIfUnset(data interface{})
	Generate(cfg *Config) (Document, error)
}

type File interface {
	Generate(cfg *Config) (Document, error)
}

type Option interface {
	apply(cfg *Config) error
}

type optionFunc func(*Config) error

func (o optionFunc) apply(s *Config) error {
	if err := o(s); err != nil {
		return err
	}
	return nil
}

func SetDescription(desc string) Option {
	return optionFunc(func(cfg *Config) error {
		cfg.Description = desc
		return nil
	})
}
func SkipChildren(skip bool) Option {
	return optionFunc(func(cfg *Config) error {
		cfg.SkipChildren = skip
		return nil
	})
}

var ErrTemplateExtensionEmpty = errors.New("Template extension can not be empty")

func SetTemplateExtension(extension string) Option {
	return optionFunc(func(cfg *Config) error {
		if extension == "" {
			return ErrTemplateExtensionEmpty
		}
		cfg.TemplateExtensionSuffix = extension
		return nil
	})
}

func Generate(file File, opts ...Option) error {

	cfg := DefaultConfig

	for _, opt := range opts {
		if err := opt.apply(&cfg); err != nil {
			return fmt.Errorf("unable to apply option %w", err)
		}
	}

	if err := generate(file, &cfg); err != nil {
		return err
	}

	return nil
}

func generate(file File, cfg *Config) error {
	doc, err := file.Generate(cfg)
	if err != nil {
		return err
	}
	if docs, ok := isDocs(doc); ok {
		for _, doc = range docs {
			if err := Write(doc); err != nil {
				return err
			}
		}
	} else {
		if err := Write(doc); err != nil {
			return err
		}
	}

	return nil
}

func isDocs(doc Document) ([]Document, bool) {
	if docs, ok := doc.(Docs); ok {
		return docs, true
	}
	return nil, false
}

func Write(file Document) error {
	if w, ok := file.(DocWriter); ok {
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
