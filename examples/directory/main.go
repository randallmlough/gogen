package main

import (
	"github.com/randallmlough/gogen"
	"log"
)

func main() {
	type Data struct {
		Name           string
		Location       string
		Age            *int
		FavoriteThings []string
		Greeting       string
		PrimaryKey     interface{}
	}
	data := Data{
		Name:           "john snow",
		Location:       "the wall",
		Age:            nil,
		FavoriteThings: []string{"dogs", "cold places", "sam"},
		Greeting:       "Hello to you too.",
		PrimaryKey:     int(1),
	}
	dir := &gogen.Dir{
		OutputDir:    "examples/directory/output",
		TemplateDir:  "examples/directory/templates",
		TemplateData: data,
	}
	if err := gogen.Generate(dir, gogen.SkipChildren(false)); err != nil {
		log.Fatal(err)
	}
}
