package main

import (
	"github.com/randallmlough/gogen"
	"log"
)

func main() {
	if err := generateGoCode(); err != nil {
		log.Fatal(err)
	}
	if err := generateFile(); err != nil {
		log.Fatal(err)
	}
}

func generateGoCode() error {
	data := map[string]interface{}{
		"Name":     "john",
		"Greeting": "Hello to you too.",
	}
	contents, err := gogen.LoadTemplate("examples/simple/gocode.gotpl")
	if err != nil {
		return err
	}
	gocode := &gogen.Go{
		Template:     contents,
		Filename:     "examples/simple/output/hello.go",
		PackageName:  "testing",
		TemplateData: data,
	}

	if err := gogen.Generate(gocode, &gogen.Config{
		GeneratedHeader: true,
		Description:     "// file generation is awesome"}); err != nil {
		return err
	}
	return nil
}
func generateFile() error {
	type Data struct {
		Name           string
		Location       string
		Age            *int
		FavoriteThings []string
	}
	data := Data{
		Name:           "john snow",
		Location:       "the wall",
		Age:            nil,
		FavoriteThings: []string{"dogs", "cold places", "sam"},
	}
	contents, err := gogen.LoadTemplate("examples/simple/file.gotpl")
	if err != nil {
		return err
	}
	doc := &gogen.Document{
		Template:     contents,
		Filename:     "examples/simple/output/file.yml",
		TemplateData: data,
	}
	if err := gogen.Generate(doc); err != nil {
		return err
	}
	return nil
}
