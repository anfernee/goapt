package chesksum

import "testing"

func TestVerify(t *testing.T) {
	tests := []struct {
		desc      string
		checksum  Checksum
		path      string
		expect    bool
		expectErr bool
	}{
		{
			desc: "md5",
			checksum: Checksum{
				Type:  MD5,
				Value: "fac2ef43e8a533db006395fb2f5769e3",
			},
			path:   "testdata/random.txt",
			expect: true,
		},
		{
			desc: "bad md5",
			checksum: Checksum{
				Type:  MD5,
				Value: "xxxxx",
			},
			path:   "testdata/random.txt",
			expect: false,
		},
		{
			desc: "sha1",
			checksum: Checksum{
				Type:  SHA1,
				Value: "d68fb8e2f5cda5e563e44d50644ebfd410d9c276",
			},
			path:   "testdata/random.txt",
			expect: true,
		},
		{
			desc: "bad sha1",
			checksum: Checksum{
				Type:  SHA1,
				Value: "xxxxx",
			},
			path:   "testdata/random.txt",
			expect: false,
		},
		{
			desc:      "wrong path",
			path:      "testdata/not-exist",
			expectErr: true,
		},
	}

	for _, test := range tests {
		got, err := Verify(test.path, test.checksum)

		if test.expectErr && err == nil {
			t.Errorf("%v: expect error; got nil error", test.desc)
		} else if !test.expectErr && err != nil {
			t.Errorf("%v: expect nil error; got %v", test.desc, err)
		}

		if got != test.expect {
			t.Errorf("%v: expect %v; got %v", test.desc, test.expect, got)
		}
	}
}
