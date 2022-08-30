package common

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ReaderOf loads io.ReadCloser from a path or url
func ReaderOf(pathOrUrl string) (io.ReadCloser, error) {
	if !strings.HasPrefix(pathOrUrl, "http") {
		return os.Open(pathOrUrl)
	}

	resp, err := http.Get(pathOrUrl)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("failed to fetch %q, status: %v", pathOrUrl, resp.Status)
	}

	return resp.Body, nil
}
