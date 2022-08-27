package release

import (
	"bufio"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Checksum string

var (
	md5  Checksum = "MD5Sum"
	sha1 Checksum = "SHA1"
)

// Release is a deb type source entry hosted via web server.
// Check sources.list(5) for details.
type Release struct {
	Origin      string
	Label       string
	Suite       string
	Version     string
	Codename    string
	Date        time.Time
	Archs       []string
	Components  []string
	Description string
	Files       map[string]*File
}

// File is a single file in a deb release.
type File struct {
	Name string
	Size int
	MD5  string
	SHA1 string
}

// Load loads a release from a url.
func Load(url string) (*Release, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parse(resp.Body)
}

func parse(r io.Reader) (*Release, error) {
	var (
		b        = bufio.NewReader(r)
		checksum Checksum
	)

	release := &Release{
		Files: map[string]*File{},
	}

	for {
		lb, _, err := b.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		line := strings.TrimSpace(string(lb))

		// Checksum Section
		if checksum != "" {
			splitter := regexp.MustCompile(" +")
			columns := splitter.Split(line, -1)
			addOrUpdate(columns, release.Files, checksum)
		}

		trim := func(line, prefix string) string {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}

		// Metadata section
		switch {
		case strings.HasPrefix(line, "Origin:"):
			release.Origin = trim(line, "Origin:")
		case strings.HasPrefix(line, "Label:"):
			release.Label = trim(line, "Label:")
		case strings.HasPrefix(line, "Suite:"):
			release.Suite = trim(line, "Suite:")
		case strings.HasPrefix(line, "Version:"):
			release.Version = trim(line, "Version:")
		case strings.HasPrefix(line, "Codename:"):
			release.Codename = trim(line, "Codename:")
		case strings.HasPrefix(line, "Description:"):
			release.Description = trim(line, "Description:") // Is it possible to have multi-line description?
		case strings.HasPrefix(line, "Architectures:"):
			release.Archs = strings.Split(trim(line, "Architectures:"), " ")
		case strings.HasPrefix(line, "Components:"):
			release.Components = strings.Split(trim(line, "Components:"), " ")
		case strings.HasPrefix(line, "Date:"):
			release.Date, _ = time.Parse(time.RFC1123, trim(line, "Date:"))
		case line == "MD5Sum:":
			checksum = md5
		case line == "SHA1:":
			checksum = sha1
		}

	}

	return release, nil
}

// addOrUpdate adds or updates a file entry from a line in release file.
//
// Example:
//  782b71b245386a0f5d99a4461ddbd571        603064135 Contents-riscv64
func addOrUpdate(columes []string, files map[string]*File, checksum Checksum) {
	if len(columes) != 3 {
		return
	}

	size, _ := strconv.Atoi(columes[1])

	if _, ok := files[columes[2]]; !ok {
		files[columes[2]] = &File{
			Name: columes[2],
			Size: size,
		}
	}

	file := files[columes[2]]
	switch checksum {
	case md5:
		file.MD5 = columes[0]
	case sha1:
		file.SHA1 = columes[0]
	}
}
