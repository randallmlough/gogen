package gogen

import (
	"bytes"
	"fmt"
	"github.com/randallmlough/gogen/gocode"
	"github.com/randallmlough/gogen/template"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"github.com/randallmlough/gogen/gocode/imports"
)

type Go struct {
	// Template is a string of the entire template that
	// will be parsed and rendered. If it's empty,
	// the plugin processor will look for .gotpl files
	// in the same directory of where you wrote the plugin.
	Template string
	Bundle   string
	// Filename is the name of the file that will be
	// written to the system disk once the template is rendered.
	Filename string
	Funcs    template.FuncMap
	// PackageName is a helper that specifies the package header declaration.
	// In other words, when you write the template you don't need to specify `package X`
	// at the top of the file. By providing PackageName in the Options, the Generate
	// function will do that for you.
	PackageName string
	Imports     []string

	packages *gocode.Packages

	TemplateData interface{}
	data         Data
}

func (g *Go) init() {
	if g.packages == nil {
		g.packages = &gocode.Packages{}
	}

	// prefetch all packages in one big packages.Load call
	g.packages.LoadAll(g.Imports...)

	gocode.CurrentImports = &gocode.Imports{Packages: g.packages, DestDir: filepath.Dir(g.Filename)}

	funcs := gocode.Funcs()
	for n, f := range g.Funcs {
		funcs[n] = f
	}
	g.Funcs = funcs

}

func (g *Go) Path() string {
	return g.Filename
}
func (g *Go) Bytes() []byte {
	return g.data.Bytes()
}
func (g *Go) Generate(cfg *Config) (File, error) {
	if gocode.CurrentImports != nil {
		panic(fmt.Errorf("recursive or concurrent call to RenderToFile detected"))
	}

	g.init()

	var buf *bytes.Buffer
	var err error
	if g.Template != "" {
		t := template.New("").Funcs(g.Funcs)
		var err error
		t, err = t.Add("goTemplate.gotpl").Parse(g.Template)
		if err != nil {
			return nil, errors.Wrap(err, "error with provided template")
		}

		buf, err = t.Execute(g.Filename, g.TemplateData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create go file")
		}
	} else {
		// if bundles is empty use location of calling function for template lookup
		rootDir := g.Bundle
		if rootDir == "" {
			_, callerFile, _, _ := runtime.Caller(1)
			rootDir = filepath.Dir(callerFile)
		}
		t := template.New("").Funcs(g.Funcs)

		var err error
		t, err = t.GatherBundles(rootDir, true)
		if err != nil {
			return nil, errors.Wrap(err, "failed to gather go bundle")
		}

		buf, err = t.ExecuteBundles(g.TemplateData, cfg.RegionTags)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build go bundle")
		}
	}

	data := g.fileTemplate(cfg)
	_, err = buf.WriteTo(&data)
	if err != nil {
		return nil, err
	}

	gocode.CurrentImports = nil
	g.data = &data
	return g, nil
}

func (g *Go) fileTemplate(cfg *Config) bytes.Buffer {
	var header bytes.Buffer
	if cfg.GeneratedHeader {
		header.WriteString(cfg.generatedText)
	}
	if cfg.Description != "" {
		header.WriteString(cfg.Description + "\n")
	}
	header.WriteString("package ")

	pkgName := g.PackageName
	if pkgName == "" {
		pkgName = gocode.PackageNameFromFile(g.Filename)
	}
	header.WriteString(pkgName)
	header.WriteString("\n\n")
	if cfg.FileNotice {
		header.WriteString(cfg.FileNoticeText)
		header.WriteString("\n\n")
	}
	header.WriteString("import (\n")
	header.WriteString(gocode.CurrentImports.String())
	header.WriteString(")\n")
	return header
}

func (g *Go) Write(file File) error {
	if err := makeDir(filepath.Dir(file.Path())); err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	formatted, err := imports.Prune(file.Path(), file.Bytes(), g.packages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gofmt failed on %s: %s\n", filepath.Base(file.Path()), err.Error())
		formatted = file.Bytes()
	}
	if err := createFile(file.Path(), formatted); err != nil {
		return errors.Wrapf(err, "failed to write %s", file.Path())
	}

	g.packages.Evict(gocode.ImportPathForDir(filepath.Dir(g.Filename)))

	return nil
}
