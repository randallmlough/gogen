package gocode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPkgAndType(t *testing.T) {
	tests := []struct {
		path    string
		Package string
		Type    string
	}{
		{path: "github.com/package/blah.Type", Package: "github.com/package/blah", Type: "Type"},
	}
	t.Run("separate package and type", func(t *testing.T) {
		for _, tt := range tests {
			Package, Type := PkgAndType(tt.path)
			require.Equal(t, tt.Package, Package)
			require.Equal(t, tt.Type, Type)
		}
	})
}

func TestNormalizeVendor(t *testing.T) {
	require.Equal(t, "bar/baz", NormalizeVendor("foo/vendor/bar/baz"))
	require.Equal(t, "[]bar/baz", NormalizeVendor("[]foo/vendor/bar/baz"))
	require.Equal(t, "*bar/baz", NormalizeVendor("*foo/vendor/bar/baz"))
	require.Equal(t, "*[]*bar/baz", NormalizeVendor("*[]*foo/vendor/bar/baz"))
	require.Equal(t, "[]*bar/baz", NormalizeVendor("[]*foo/vendor/bar/baz"))
}

func TestSanitizePackageName(t *testing.T) {
	require.Equal(t, "foo", SanitizePackageName("foo"))
	require.Equal(t, "baz", SanitizePackageName("foo/bar/baz"))
	require.Equal(t, "foo_bar", SanitizePackageName("foo-bar"))
	require.Equal(t, "bar", SanitizePackageName("foo/bar"))
}

func TestPackageNameFromFile(t *testing.T) {
	require.Equal(t, "foo", PackageNameFromFile("foo"))
	require.Equal(t, "foo", PackageNameFromFile("/foo"))
	require.Equal(t, "foo", PackageNameFromFile("./foo"))
	require.Equal(t, "foo", PackageNameFromFile("/foo/"))
	require.Equal(t, "foo", PackageNameFromFile("foo/"))
	require.Equal(t, "foo", PackageNameFromFile("foo.go"))
	require.Equal(t, "foo", PackageNameFromFile("/foo.go"))
	require.Equal(t, "foo", PackageNameFromFile("./foo.go"))
	require.Equal(t, "package", PackageNameFromFile("package/foo.go"))
	require.Equal(t, "package", PackageNameFromFile("/package/foo.go"))
	require.Equal(t, "package", PackageNameFromFile("./package/foo.go"))
	require.Equal(t, "package", PackageNameFromFile("parent/package/foo.go"))
	require.Equal(t, "package", PackageNameFromFile("/parent/package/foo.go"))
	require.Equal(t, "package", PackageNameFromFile("./parent/package/foo.go"))
}
