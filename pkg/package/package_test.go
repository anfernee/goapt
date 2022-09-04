package pkg

import (
	"testing"
)

func TestLoadPackage(t *testing.T) {
	paths := []string{
		"testdata/bazel-packages.gz",
		"testdata/bazel-packages.xz",
	}
	for _, path := range paths {
		pkgs, err := Load(path)
		if err != nil {
			t.Error(err)
		}
		if len(pkgs) != 66 {
			t.Errorf("expect len(pkgz)==66; got %d", len(pkgs))
		}

		for _, pkg := range pkgs {
			if pkg.Filename == "" {
				t.Errorf("expect non empty pkg.Filename; got <empty>")
			}
			if pkg.Name == "" {
				t.Errorf("expect non empty pkg.Name; got <empty>")
			}
			if pkg.Arch == "" {
				t.Errorf("expect non empty pkg.Arch; got <empty>")
			}
			if pkg.Size == 0 {
				t.Errorf("expect pkg.Size > 0; got 0")
			}
		}
	}
}
