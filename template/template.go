package template

import (
	"bytes"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

func New(name string) *Template {
	return &Template{
		template: template.New(name).Funcs(StandardTemplateFuncs),
		bundles:  nil,
	}
}

type Template struct {
	template *template.Template
	bundles  []string
}

func (t *Template) clone() *Template {
	if t == nil {
		return New("")
	}
	var c Template
	c = *t
	return &c
}

func (t *Template) Funcs(funcs template.FuncMap) *Template {
	return &Template{template: t.template.Funcs(funcs), bundles: t.bundles}
}

func (t *Template) Add(name string) *Template {
	return &Template{template: t.template.New(name), bundles: t.bundles}
}

func (t *Template) Parse(contents string) (*Template, error) {
	tmp, err := t.template.Parse(contents)
	if err != nil {
		return nil, err
	}
	return &Template{template: tmp, bundles: t.bundles}, nil
}

func (t *Template) ParseTemplate(pathToTemplate string) (*Template, error) {
	b, err := ioutil.ReadFile(pathToTemplate)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read template %q", pathToTemplate)
	}
	return t.Parse(string(b))
}

func (t *Template) GatherBundles(rootDir, templateExtension string, skipChildren bool) (*Template, error) {
	tt := t.clone()
	bundles := []string{}
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip child directories
		if info.IsDir() && (rootDir != path && skipChildren) {
			return filepath.SkipDir
		}
		templateName := filepath.ToSlash(strings.TrimPrefix(path, rootDir+string(os.PathSeparator)))
		if !strings.HasSuffix(info.Name(), templateExtension) {
			return nil
		}

		tt, err = t.Add(templateName).ParseTemplate(path)
		if err != nil {
			return errors.Wrap(err, path)
		}

		bundles = append(bundles, templateName)

		return nil
	})
	if err != nil {
		return tt, errors.Wrap(err, "failed to gather bundles")
	}
	sortTemplates(bundles, templateExtension)
	tt.bundles = bundles
	return tt, nil
}

func (t *Template) Execute(filename string, data interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := t.template.Execute(&buf, data)
	if err != nil {
		return nil, errors.Wrap(err, filename)
	}
	return &buf, nil
}

func (t *Template) ExecutePart(buf *bytes.Buffer, part string, data interface{}) error {
	err := t.template.Lookup(part).Execute(buf, data)
	if err != nil {
		return errors.Wrap(err, part)
	}
	return nil
}

func (t *Template) ExecuteBundles(data interface{}, addRegionTag bool) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	if t.bundles == nil {
		if err := t.template.Execute(&buf, data); err != nil {
			return nil, err
		}
	} else {
		for _, bundle := range t.bundles {
			if addRegionTag {
				beginRegion(buf, bundle)
			}
			err := t.ExecutePart(&buf, bundle, data)
			if err != nil {
				return nil, errors.Wrap(err, bundle)
			}
			if addRegionTag {
				endRegion(buf, bundle)
			}
		}
	}

	return &buf, nil
}

func sortTemplates(bundles []string, templateExtension string) {
	// then execute all the important looking ones in order, adding them to the same file
	sort.Slice(bundles, func(i, j int) bool {
		// important files go first
		if strings.HasSuffix(bundles[i], "!"+templateExtension) {
			return true
		}
		if strings.HasSuffix(bundles[j], "!"+templateExtension) {
			return false
		}
		return bundles[i] < bundles[j]
	})
}
