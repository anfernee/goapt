package release

import (
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestReleaseParse(t *testing.T) {
	d, _ := os.Open("testdata/example-focal.txt")
	r, err := parse(d)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	expected := &Release{
		Origin:      "Ubuntu",
		Label:       "Ubuntu",
		Suite:       "focal",
		Version:     "20.04",
		Codename:    "focal",
		Description: "Ubuntu Focal 20.04",
		Date:        time.Date(2020, time.April, 23, 17, 33, 17, 0, time.UTC),
		Archs: []string{
			"amd64",
			"arm64",
		},
		Components: []string{
			"main",
			"restricted",
			"universe",
			"multiverse",
		},
		Files: map[string]*File{
			"main/binary-amd64/Packages": {
				Name: "main/binary-amd64/Packages",
				Size: 5826751,
				MD5:  "7ef83228ec207df10acac48fbdd81112",
				SHA1: "aef5c36ce45bd5c3154a1bb03c62b6cfb33e2bc6",
			},
			"main/binary-amd64/Packages.gz": {
				Name: "main/binary-amd64/Packages.gz",
				Size: 1274738,
				MD5:  "737716846a5e245ed9e7590d482005fe",
				SHA1: "9f9fc72cf9069945408ca84e7a7e0e22d3c2dda4",
			},
		},
	}
	if !cmp.Equal(expected, r) {
		t.Errorf("unexpected diff: %v", cmp.Diff(expected, r))
	}
}
