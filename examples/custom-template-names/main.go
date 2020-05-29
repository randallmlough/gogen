package main

import (
	"github.com/randallmlough/gogen"
)

func main() {
	data := map[string]interface{}{
		"Name":     "john",
		"Greeting": "Hello to you too.",
	}
	contents, err := gogen.LoadTemplate("examples/custom-template-names/gocode.tpl")
	if err != nil {
		panic(err)
	}
	gocode := &gogen.Go{
		Template:     contents,
		Filename:     "examples/custom-template-names/output/hello.go",
		PackageName:  "testing", // you can be explicit or let us guess based on the path
		TemplateData: data,
	}

	if err := gogen.Generate(gocode,
		gogen.SetTemplateExtension(".tpl"),
		gogen.SetDescription("The possibilities of file generation are endless"),
	); err != nil {
		panic(err)
	}
}
