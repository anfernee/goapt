package pkg

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadSingleDebSource(t *testing.T) {
	expected := []string{
		"http://archive.ubuntu.com/ubuntu/dists/focal/main/binary-amd64/Packages",
		"http://archive.ubuntu.com/ubuntu/dists/focal/restricted/binary-amd64/Packages",
		"http://archive.ubuntu.com/ubuntu/dists/focal/main/source/Sources",
		"http://archive.ubuntu.com/ubuntu/dists/focal/restricted/source/Sources",
		"http://archive.ubuntu.com/ubuntu/dists/focal-updates/main/binary-amd64/Packages",
		"http://archive.ubuntu.com/ubuntu/dists/focal-updates/restricted/binary-amd64/Packages",
		"http://archive.ubuntu.com/ubuntu/dists/focal-updates/main/source/Sources",
		"http://archive.ubuntu.com/ubuntu/dists/focal-updates/restricted/source/Sources",
	}
	list, err := loadDebianSourceList("testdata/single_sources.list")
	if err != nil {
		t.Error(err)
	}

	got := []string{}
	for _, ds := range list {
		got = append(got, ds.ResourceURL())
	}

	if !cmp.Equal(expected, got) {
		t.Errorf("unexpected diff: %v", cmp.Diff(expected, got))
	}
}

func TestLoadDebSource(t *testing.T) {
	expected := []string{
		"http://archive.ubuntu.com/ubuntu/dists/focal/main/binary-amd64/Packages",
		"http://archive.ubuntu.com/ubuntu/dists/focal/restricted/binary-amd64/Packages",
		"http://archive.ubuntu.com/ubuntu/dists/focal/main/source/Sources",
		"http://archive.ubuntu.com/ubuntu/dists/focal/restricted/source/Sources",
		"http://archive.ubuntu.com/ubuntu/dists/focal-updates/main/binary-amd64/Packages",
		"http://archive.ubuntu.com/ubuntu/dists/focal-updates/restricted/binary-amd64/Packages",
		"http://archive.ubuntu.com/ubuntu/dists/focal-updates/main/source/Sources",
		"http://archive.ubuntu.com/ubuntu/dists/focal-updates/restricted/source/Sources",
	}
	list, err := loadDebianSourceList("testdata/sources.list")
	if err != nil {
		t.Error(err)
	}

	got := []string{}
	for _, ds := range list {
		got = append(got, ds.ResourceURL())
	}

	if !cmp.Equal(expected, got) {
		t.Errorf("unexpected diff: %v", cmp.Diff(expected, got))
	}
}
