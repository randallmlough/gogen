package gogen

import (
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func MustLoadTemplate(pathToTemplate string) string {
	contents, err := LoadTemplate(pathToTemplate)
	if err != nil {
		panic(err)
	}
	return contents
}
func LoadTemplate(pathToTemplate string) (string, error) {

	b, err := ioutil.ReadFile(pathToTemplate)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func findGoNamedType(def types.Type) (*types.Named, error) {
	if def == nil {
		return nil, nil
	}

	namedType, ok := def.(*types.Named)
	if !ok {
		return nil, errors.Errorf("expected %s to be a named type, instead found %T\n", def.String(), def)
	}

	return namedType, nil
}

func findGoInterface(def types.Type) (*types.Interface, error) {
	if def == nil {
		return nil, nil
	}
	namedType, err := findGoNamedType(def)
	if err != nil {
		return nil, err
	}
	if namedType == nil {
		return nil, nil
	}

	underlying, ok := namedType.Underlying().(*types.Interface)
	if !ok {
		return nil, errors.Errorf("expected %s to be a named interface, instead found %s", def.String(), namedType.String())
	}

	return underlying, nil
}

func equalFieldName(source, target string) bool {
	source = strings.Replace(source, "_", "", -1)
	target = strings.Replace(target, "_", "", -1)
	return strings.EqualFold(source, target)
}

// getWorkingPath gets the current working directory
func getWorkingPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return wd, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func makeDirForFile(fileName string) error {
	dir := filepath.Dir(fileName)
	if err := makeDir(dir); err != nil {
		return fmt.Errorf("unable to create dir for file: %s %w", fileName, err)
	}
	return nil
}

func makeDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("unable to create dir: %s %w", dir, err)
	}
	return nil
}

func createFile(fileName string, data []byte) error {
	if err := ioutil.WriteFile(fileName, data, 0644); err != nil {
		return fmt.Errorf("unable to create file: %s %w", fileName, err)
	}
	return nil
}
