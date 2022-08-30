package release

import (
	"bufio"
	"os"
	"testing"
)

func TestVerifyWithOptions(t *testing.T) {
	tests := []struct {
		desc          string
		options       *VerifyOptions
		cleartextPath string
		expectErr     bool
	}{
		{
			desc:          "local cleartext, unarmored key",
			cleartextPath: "testdata/inrelease.txt",
			options: &VerifyOptions{
				KeyPath: "testdata/bazel-archive-keyring.gpg",
				Armored: false,
			},
		},
		{
			desc:          "local cleartext, armored key",
			cleartextPath: "testdata/inrelease.txt",
			options: &VerifyOptions{
				KeyPath: "testdata/bazel-release.pub.gpg",
				Armored: true,
			},
		},
	}

	for _, test := range tests {
		_, err := VerifyWithOptions(test.cleartextPath, test.options)
		if test.expectErr && err == nil {
			t.Errorf("expect err; got nil")
		} else if !test.expectErr && err != nil {
			t.Errorf("expect nil err; got %q", err)
		}
	}
}

func TestUbuntuKnwonPath(t *testing.T) {
	if !isUbuntu() {
		t.Skip("the test only works in ubuntu")
	}

	url := "http://us.archive.ubuntu.com/ubuntu/dists/focal/InRelease"
	_, err := Verify(url)
	if err != nil {
		t.Errorf("failed to verify %s on ubuntu: %q.\nTry to run 'apt update' to see if it works.", url, err)
	}
}

func isUbuntu() bool {
	f, err := os.Open("/etc/lsb-release")
	if err != nil {
		return false
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			return false
		}
		if string(line) == "DISTRIB_ID=Ubuntu" {
			return true
		}
	}
}
