package chesksum

import (
	"crypto"
	"fmt"
	"io"

	"github.com/anfernee/goapt/pkg/common"
)

type Type string

var (
	MD5  Type = "MD5Sum"
	SHA1 Type = "SHA1"

	hashMap = map[Type]crypto.Hash{
		MD5:  crypto.MD5,
		SHA1: crypto.SHA1,
	}
)

type Checksum struct {
	Type  Type
	Value string
}

// Verify checksum of a file locally or hosted on http server.
func Verify(pathOrUrl string, checksum Checksum) (bool, error) {
	hash, ok := hashMap[checksum.Type]
	if !ok {
		return false, fmt.Errorf("unsupported checksum type %v", checksum.Type)
	}

	rc, err := common.ReaderOf(pathOrUrl)
	if err != nil {
		return false, err
	}
	defer rc.Close()

	h := hash.New()
	if _, err := io.Copy(h, rc); err != nil {
		return false, err
	}

	fmt.Printf("%x\n", h.Sum(nil))

	return fmt.Sprintf("%x", h.Sum(nil)) == checksum.Value, nil
}
