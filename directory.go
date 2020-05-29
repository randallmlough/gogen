package gogen

import (
	"github.com/pkg/errors"
	"github.com/randallmlough/gogen/gocode"
	"os"
	"path/filepath"
	"strings"
)

type Directory struct {
	OutputDir    string
	TemplateDir  string
	TemplateData interface{}
}

var _ Gen = (*Directory)(nil)

func (d *Directory) Name() string {
	return d.OutputDir
}

func (d *Directory) SetTemplateDataIfUnset(data interface{}) {
	if d.TemplateData == nil {
		d.TemplateData = data
	}
}

// walk directories and create templates from files
func (d *Directory) Generate(cfg *Config) (File, error) {

	if err := cfg.check(); err != nil {
		return nil, errors.Wrap(err, "config is improperly formatted")
	}

	ff, err := d.Files(cfg)
	if err != nil {
		return nil, err
	}
	files := Files{}
	for _, f := range ff {
		if d, err := f.Generate(cfg); err != nil {
			return nil, err
		} else {
			files = append(files, d)
		}
	}
	return files, nil
}

type fileable interface {
	Generate(cfg *Config) (File, error)
	Path() string
}

func (d *Directory) Files(cfg *Config) ([]fileable, error) {
	rootDir := d.TemplateDir

	files := []fileable{}
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if rootDir != path && (info.IsDir() && cfg.SkipChildren) {
			return filepath.SkipDir
		}

		if !strings.HasSuffix(info.Name(), cfg.TemplateExtensionSuffix) {
			return nil
		}

		contents, err := LoadTemplate(path)
		if err != nil {
			return err
		}
		path = strings.TrimSuffix(path, cfg.TemplateExtensionSuffix)

		if d.OutputDir != "" {
			path = strings.Replace(path, rootDir, d.OutputDir, 1)
		}

		switch filepath.Ext(path) {
		case ".go":
			code := Go{
				Template:     contents,
				Filename:     path,
				PackageName:  gocode.PackageNameFromFile(path),
				TemplateData: d.TemplateData,
			}
			files = append(files, &code)
		default:
			file := Document{
				Template:     contents,
				Filename:     path,
				TemplateData: d.TemplateData,
			}
			files = append(files, &file)
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "locating templates")
	}
	return files, nil
}

// look for a base template file
