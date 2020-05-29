package gogen

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/randallmlough/gogen/template"
	"path/filepath"
	"runtime"
)

type Document interface {
	Path() string
	Data
}

type Data interface {
	Bytes() []byte
}

type DocWriter interface {
	Write(file Document) error
}

type Docs []Document

// should never be called. Used to implement Document interface
func (d Docs) Bytes() []byte {
	return nil
}

// should never be called. Used to implement Document interface
func (d Docs) String() string {
	return ""
}

// should never be called. Used to implement Document interface
func (d Docs) Path() string {
	return ""
}

type Doc struct {
	// Template is a string of the entire template that
	// will be parsed and rendered. If it's empty,
	// the plugin processor will look for .gotpl files
	// in the same directory of where you wrote the plugin.
	Template string
	Bundle   string
	// Filename is the name of the file that will be
	// written to the system disk once the template is rendered.
	Filename string

	Funcs template.FuncMap

	TemplateData interface{}
	data         Data
}

func (doc *Doc) Path() string {
	return doc.Filename
}

func (doc *Doc) Bytes() []byte {
	return doc.data.Bytes()
}

func (doc *Doc) SetTemplateDataIfUnset(data interface{}) {
	if doc.TemplateData == nil {
		doc.TemplateData = data
	}
}

func (doc *Doc) Generate(cfg *Config) (Document, error) {

	if err := cfg.check(); err != nil {
		return nil, errors.Wrap(err, "config is improperly formatted")
	}

	var buf *bytes.Buffer
	var err error
	if doc.Template != "" {
		t := template.New("").Funcs(doc.Funcs)
		var err error
		t, err = t.Add("fileTemplate" + cfg.TemplateExtensionSuffix).Parse(doc.Template)
		if err != nil {
			return nil, errors.Wrap(err, "error with provided template")
		}

		buf, err = t.Execute(doc.Filename, doc.TemplateData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate single go file")
		}
	} else {
		rootDir := doc.Bundle
		if rootDir == "" {
			_, callerFile, _, _ := runtime.Caller(1)
			rootDir = filepath.Dir(callerFile)
		}

		t := template.New("").Funcs(doc.Funcs)

		var err error
		t, err = t.GatherBundles(rootDir, cfg.TemplateExtensionSuffix, cfg.SkipChildren)
		if err != nil {
			return nil, errors.Wrap(err, "failed to gather document bundle")
		}

		buf, err = t.ExecuteBundles(doc.TemplateData, cfg.RegionTags)
		if err != nil {
			return nil, errors.Wrap(err, "failed to build document bundle")
		}
	}
	data := doc.fileTemplate(cfg)
	_, err = buf.WriteTo(&data)
	if err != nil {
		return nil, err
	}

	doc.data = &data
	return doc, nil
}

func (doc *Doc) fileTemplate(cfg *Config) bytes.Buffer {
	var header bytes.Buffer
	if cfg.GeneratedHeader {
		header.WriteString(cfg.generatedText)
	}
	if cfg.Description != "" {
		header.WriteString(cfg.Description + "\n")
	}
	if cfg.FileNotice {
		header.WriteString(cfg.FileNoticeText)
		header.WriteString("\n\n")
	}
	return header
}
