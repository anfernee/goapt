package release

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/ProtonMail/gopenpgp/v2/helper"
)

const (
	defaultTrustedPath = "/etc/apt/trusted.gpg"
	defaultTrustedDir  = "/etc/apt/trusted.gpg.d"
)

// VerifyOptions specifies the option to verify cleartext GPG
// signature.
type VerifyOptions struct {
	KeyPath      string
	Armored      bool
	AutoDiscover bool
}

// VerifyWithOptions verifies a local file or http/https url with given public GNG
// public key.
func VerifyWithOptions(path string, options *VerifyOptions) (string, error) {
	cleartext, err := loadClearText(path)
	if err != nil {
		return "", err
	}

	if options == nil || options.AutoDiscover {
		return verifyWithKnownKeys(cleartext)
	}

	return verifyWithKeyRing(cleartext, options.KeyPath, options.Armored)
}

// Verify verifies a local file or http/https url with well known public GNG
// keys saved by apt-key, in /etc/apt/trusted.gpg and under /etc/apt/trusted.gpg.d
func Verify(path string) (string, error) {
	return VerifyWithOptions(path, nil)
}

// loadClearText loads cleartext message from path or url.
func loadClearText(pathOrUrl string) ([]byte, error) {
	if !strings.HasPrefix(pathOrUrl, "http") {
		return ioutil.ReadFile(pathOrUrl)
	}

	resp, err := http.Get(pathOrUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %q, status: %v", pathOrUrl, resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}

// loadKeyRing loads keyring from a key path. armored specifies whether the key file
// is armored or not.
func loadKeyRing(keyPath string, armored bool) (*crypto.KeyRing, error) {
	d, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	var pubkey *crypto.Key
	if armored {
		pubkey, err = crypto.NewKeyFromArmored(string(d))
	} else {
		pubkey, err = crypto.NewKey(d)
	}
	if err != nil {
		return nil, err
	}

	return crypto.NewKeyRing(pubkey)
}

// verifyWithKnownKeys verifies clear text message with known keys saved in /etc/apt/trusted.gpg
// and under /etc/apt/trusted.gpg.d
func verifyWithKnownKeys(text []byte) (string, error) {
	keyFiles := []string{defaultTrustedPath}

	entries, err := os.ReadDir(defaultTrustedDir)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".gpg") {
				keyFiles = append(keyFiles, filepath.Join(defaultTrustedDir, entry.Name()))
			}
		}
	}

	for _, keyFile := range keyFiles {
		if s, err := verifyWithKeyRing(text, keyFile, false); err == nil {
			return s, nil
		}
	}

	return "", fmt.Errorf("failed to verify cleartext message")
}

// verifyWithKeyRing verifies clear text message with a keyring specified in keyFile.
func verifyWithKeyRing(text []byte, keyFile string, armored bool) (string, error) {
	keyRing, err := loadKeyRing(keyFile, armored)
	if err != nil {
		return "", err
	}

	return helper.VerifyCleartextMessage(keyRing, string(text), crypto.GetUnixTime())
}
