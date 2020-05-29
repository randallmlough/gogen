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
		"PrimaryKey":        1,
		"PartTwoAnswer":     2,
		"PartThreeCtxValue": "some-random-key",
	}
	gocode := &gogen.Go{
		Bundle:       "examples/bundles/templates",
		Filename:     "examples/bundles/output/bundle.go",
		PackageName:  "testing",
		TemplateData: data,
	}
	if err := gogen.Generate(gocode); err != nil {
		return err
	}
	return nil
}

func generateFile() error {
	data := map[string]interface{}{
		"PartOne":        "i'm a bundle!",
		"FavoriteThings": []string{"organization", "fun code"},
		"PartTwo":        "oh the possibilities",
	}
	gocode := &gogen.Document{
		Bundle:       "examples/bundles/templates/file",
		Filename:     "examples/bundles/output/bundle.yml",
		TemplateData: data,
	}
	if err := gogen.Generate(gocode); err != nil {
		return err
	}
	return nil
}
