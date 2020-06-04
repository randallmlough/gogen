package imports

import (
	"io/ioutil"
	"testing"

	"github.com/randallmlough/gogen/gocode"
	"github.com/stretchr/testify/require"
)

func TestPrune(t *testing.T) {
	// prime the packages cache so that it's not considered uninitialized

	b, err := Prune("testdata/unused.go", mustReadFile("testdata/unused.go"), &gocode.Packages{})
	require.NoError(t, err)
	require.Equal(t, string(mustReadFile("testdata/unused.expected.go")), string(b))
}

func mustReadFile(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return b
}
