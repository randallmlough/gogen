package gogen

import (
	"testing"
)

func TestDirectory_Files(t *testing.T) {

	t.Run("return all the .gotpl files", func(t *testing.T) {

		dir := Dir{
			OutputDir:   "",
			TemplateDir: "testdata/walk",
		}
		cfg := &Config{
			SkipChildren:            false,
			TemplateExtensionSuffix: defaultTemplateExtension,
		}
		files, err := dir.Files(cfg)
		if err != nil {
			t.Errorf("walkDirectory() error = %v, wantErr %v", err, false)
		}
		wantFiles := 4
		if len(files) != wantFiles {
			t.Errorf("walkDirectory() got %d files, wanted %d", len(files), wantFiles)
		}
	})

	t.Run("return only the .gotpl files in the rootdir", func(t *testing.T) {

		dir := Dir{
			OutputDir:   "",
			TemplateDir: "testdata/walk",
		}
		cfg := &Config{
			SkipChildren:            true,
			TemplateExtensionSuffix: defaultTemplateExtension,
		}
		files, err := dir.Files(cfg)
		if err != nil {
			t.Errorf("walkDirectory() error = %v, wantErr %v", err, false)
		}
		wantFiles := 2
		if len(files) != wantFiles {
			t.Errorf("walkDirectory() got %d files, wanted %d", len(files), wantFiles)
		}
	})
}
