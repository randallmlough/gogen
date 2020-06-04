package gocode

import (
	"go/types"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImports(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)

	aBar := "github.com/randallmlough/gogen/gocode/testdata/a/bar"
	bBar := "github.com/randallmlough/gogen/gocode/testdata/b/bar"
	mismatch := "github.com/randallmlough/gogen/gocode/testdata/pkg_mismatch"

	t.Run("multiple lookups is ok", func(t *testing.T) {
		a := Imports{DestDir: wd, Packages: &Packages{}}

		require.Equal(t, "bar", a.Lookup(aBar))
		require.Equal(t, "bar", a.Lookup(aBar))
	})

	t.Run("lookup by type", func(t *testing.T) {
		a := Imports{DestDir: wd, Packages: &Packages{}}

		pkg := types.NewPackage("github.com/randallmlough/gogen/templates/testdata/b/bar", "bar")
		typ := types.NewNamed(types.NewTypeName(0, pkg, "Boolean", types.Typ[types.Bool]), types.Typ[types.Bool], nil)

		require.Equal(t, "bar.Boolean", a.LookupType(typ))
	})

	t.Run("duplicates are decollisioned", func(t *testing.T) {
		a := Imports{DestDir: wd, Packages: &Packages{}}

		require.Equal(t, "bar", a.Lookup(aBar))
		require.Equal(t, "bar1", a.Lookup(bBar))

		t.Run("additionial calls get decollisioned name", func(t *testing.T) {
			require.Equal(t, "bar1", a.Lookup(bBar))
		})
	})

	t.Run("package name defined in code will be used", func(t *testing.T) {
		a := Imports{DestDir: wd, Packages: &Packages{}}

		require.Equal(t, "turtles", a.Lookup(mismatch))
	})

	t.Run("string printing for import block", func(t *testing.T) {
		a := Imports{DestDir: wd, Packages: &Packages{}}
		a.Lookup(aBar)
		a.Lookup(bBar)
		a.Lookup(mismatch)

		require.Equal(
			t,
			`"github.com/randallmlough/gogen/gocode/testdata/a/bar"
bar1 "github.com/randallmlough/gogen/gocode/testdata/b/bar"
turtles "github.com/randallmlough/gogen/gocode/testdata/pkg_mismatch"`,
			a.String(),
		)
	})

	t.Run("aliased imports will not collide", func(t *testing.T) {
		a := Imports{DestDir: wd, Packages: &Packages{}}

		_, _ = a.Reserve(aBar, "abar")
		_, _ = a.Reserve(bBar, "bbar")

		require.Equal(t, `abar "github.com/randallmlough/gogen/gocode/testdata/a/bar"
bbar "github.com/randallmlough/gogen/gocode/testdata/b/bar"`, a.String())
	})

}
