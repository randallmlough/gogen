package gogen

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestTemplateOverride(t *testing.T) {
	f, err := ioutil.TempFile("", "gogen")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.RemoveAll(f.Name())

	gocode := &Go{
		Template:    "hello",
		Filename:    f.Name(),
		PackageName: "testing",
	}
	if err := Generate(gocode); err != nil {
		t.Fatal(err)
	}
}
