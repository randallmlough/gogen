package gogen

import (
	"errors"
	"fmt"
)

type Gen interface {
	Generate(cfg *Config) (File, error)
}

func Generate(file Gen, opts ...Option) error {

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

func isFiles(doc File) ([]File, bool) {
	if docs, ok := doc.(Files); ok {
		return docs, true
	}
	return nil, false
}

func generate(gen Gen, cfg *Config) error {
	file, err := gen.Generate(cfg)
	if err != nil {
		return err
	}
	if files, ok := isFiles(file); ok {
		for _, file = range files {
			if err := Write(file); err != nil {
				return err
			}
		}
	} else {
		if err := Write(file); err != nil {
			return err
		}
	}

	return nil
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
