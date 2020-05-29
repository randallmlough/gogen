package main

import "github.com/randallmlough/gogen"

type Data struct {
	Name           string
	Location       string
	Age            *int
	FavoriteThings []string
	Greeting       string
	PrimaryKey     interface{}
}

var files = []gogen.List{
	&gogen.Go{
		Filename: "examples/list/output/gocode.go",
		Template: gogen.MustLoadTemplate("examples/list/templates/gocode.go.gotpl"),
	},
	&gogen.Doc{
		Filename: "examples/list/output/file.yml",
		Template: gogen.MustLoadTemplate("examples/list/templates/file.yml.gotpl"),
	},
	&gogen.Dir{
		OutputDir:   "examples/list/output/types",
		TemplateDir: "examples/list/templates/types",
	},
}

func main() {
	data := &Data{
		Name:           "john snow",
		Location:       "the wall",
		Age:            nil,
		FavoriteThings: []string{"dogs", "cold places", "sam"},
		Greeting:       "Hello to you too.",
		PrimaryKey:     int(1),
	}
	if err := generateFiles(files, data); err != nil {
		panic(err)
	}
}

func generateFiles(files []gogen.List, data *Data) error {
	for _, file := range files {
		file.SetTemplateDataIfUnset(data)
		if err := gogen.Generate(file); err != nil {
			return err
		}
	}
	return nil
}
