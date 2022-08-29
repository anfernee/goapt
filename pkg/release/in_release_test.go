package release

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
	"github.com/google/go-cmp/cmp"
)

const expectedText = `Origin: Bazel Authors
Label: Bazel
Codename: stable
Date: Tue, 23 Aug 2022 02:01:57 UTC
Architectures: amd64
Components: jdk1.8
Description: Bazel APT Repository
MD5Sum:
 cc5d9c408a25fcd99715b80f5a7b34ef 60992 jdk1.8/binary-amd64/Packages
 4a1a6ba1e47262fc01ade424b8f00b0c 7786 jdk1.8/binary-amd64/Packages.gz
 65914fa446cb900d6ec80d4d874ca269 107 jdk1.8/binary-amd64/Release
 e90595e1bd25ddc670bad7f6a3635417 6974 jdk1.8/source/Sources.gz
 5d382ee367f605efe7cd575173ec3e42 108 jdk1.8/source/Release
SHA1:
 8e116377076927ceb11482f0857c2023f4557b33 60992 jdk1.8/binary-amd64/Packages
 e489d3605b075189355e5b0923b8b91b42b3112b 7786 jdk1.8/binary-amd64/Packages.gz
 2fd7478dde112e712451130e5523380035aba9aa 107 jdk1.8/binary-amd64/Release
 2b27996a49ec8fc40b3a36fef9010586b802f777 6974 jdk1.8/source/Sources.gz
 c8e6785823f8d1b29c67e7510878e11b2bcfffeb 108 jdk1.8/source/Release
SHA256:
 6eef9f9bd6f66654e30854dc72d255fe78f51c9716c3aeed150916c71dff4b3a 60992 jdk1.8/binary-amd64/Packages
 31ba57ba3e945288d37e8bd8173475ce29fb61da079b22efb618b50e036f0b1c 7786 jdk1.8/binary-amd64/Packages.gz
 7097d90c6b91d32fa1e16a354f73e3ad49c650d78497faf660dd7812f14ca76e 107 jdk1.8/binary-amd64/Release
 5e50c44f008092974fc27f3b409c5a7297f037777689eec25ecb674a70bc2445 6974 jdk1.8/source/Sources.gz
 c89f2c5ed362a65e74f45fe92f5079753ac4800305ed5b99963001dab0b6b0f1 108 jdk1.8/source/Release`

func TestVerify(t *testing.T) {
	/*
		pubkeyFile, _ := os.Open("testdata/bazel-release.pub.gpg")
		pubkey, err := crypto.NewKeyFromArmoredReader(pubkeyFile)
		if err != nil {
			t.Fatal(err)
		}
	*/
	d, _ := ioutil.ReadFile("testdata/bazel-archive-keyring.gpg")
	pubkey, err := crypto.NewKey(d)
	if err != nil {
		t.Fatal(err)
	}

	keyring, err := crypto.NewKeyRing(pubkey)
	if err != nil {
		t.Fatal(err)
	}

	signedMsg, _ := ioutil.ReadFile("testdata/inrelease.txt")
	fmt.Println(string(signedMsg))
	plain, err := helper.VerifyCleartextMessage(keyring, string(signedMsg), crypto.GetUnixTime())
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(plain, expectedText) {
		t.Errorf("unexpected diff between original and verified text: %v", cmp.Diff(plain, expectedText))
	}
}

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
