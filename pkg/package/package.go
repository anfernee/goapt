package pkg

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/anfernee/goapt/pkg/common"
	"github.com/ulikunitz/xz"
)

// Package is a deb package
type Package struct {
	common.Metadata
	Filename string
	Size     int
}

// Load loads packages from a path or URL
func Load(pathOrUrl string) ([]Package, error) {
	var (
		compressFmt string
		r           io.Reader
		rc          io.ReadCloser
		err         error
	)

	rc, err = readerOf(pathOrUrl)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	if strings.HasSuffix(pathOrUrl, ".xz") {
		compressFmt = "xz"
		r, err = xz.NewReader(rc)
	} else if strings.HasSuffix(pathOrUrl, ".gz") {
		compressFmt = "gzip"
		r, err = gzip.NewReader(rc)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to decompress %v", compressFmt)
	}

	return parse(r)
}

func parse(r io.Reader) ([]Package, error) {
	var (
		buf          = bufio.NewReader(r)
		ret          = []Package{Package{}}
		cur *Package = &ret[len(ret)-1]
	)
	for {
		lb, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		trim := func(line, prefix string) string {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}

		line := strings.TrimSpace(string(lb))
		switch {
		case line == "":
			ret = append(ret, Package{})
			cur = &ret[len(ret)-1]
		case strings.HasPrefix(line, "Package:"):
			cur.Name = trim(line, "Package:")
		case strings.HasPrefix(line, "Version:"):
			cur.Version = trim(line, "Version:")
		case strings.HasPrefix(line, "Section:"):
			cur.Section = trim(line, "Section:")
		case strings.HasPrefix(line, "Origin:"):
			cur.Origin = trim(line, "Origin:")
		case strings.HasPrefix(line, "Homepage:"):
			cur.Homepage = trim(line, "Homepage:")
		case strings.HasPrefix(line, "Filename:"):
			cur.Filename = trim(line, "Filename:")
		case strings.HasPrefix(line, "Size:"):
			cur.Size, _ = strconv.Atoi(trim(line, "Size:"))
		}
	}

	if cur.Name == "" {
		ret = ret[:len(ret)-1]
	}
	return ret, nil
}

// readerOf loads io.ReadCloser from a path or url
func readerOf(pathOrUrl string) (io.ReadCloser, error) {
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
