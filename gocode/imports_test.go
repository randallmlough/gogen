package gocode

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportPathForDir(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	assert.Equal(t, "github.com/randallmlough/gogen/gocode", ImportPathForDir(wd))
	assert.Equal(t, "github.com/randallmlough/gogen/api", ImportPathForDir(filepath.Join(wd, "..", "api")))

	// doesnt contain go code, but should still give a valid import path
	assert.Equal(t, "github.com/randallmlough/gogen/docs", ImportPathForDir(filepath.Join(wd, "..", "docs")))

	// directory does not exist
	assert.Equal(t, "github.com/randallmlough/dos", ImportPathForDir(filepath.Join(wd, "..", "..", "dos")))

	if runtime.GOOS == "windows" {
		assert.Equal(t, "", ImportPathForDir("C:/doesnotexist"))
	} else {
		assert.Equal(t, "", ImportPathForDir("/doesnotexist"))
	}
}

func TestNameForDir(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	fmt.Println("WD", wd)
	assert.Equal(t, "tmp", NameForDir("/tmp"))
	assert.Equal(t, "gocode", NameForDir(wd))
	assert.Equal(t, "gogen", NameForDir(wd+"/.."))
	assert.Equal(t, "randallmlough", NameForDir(wd+"/../.."))
}
