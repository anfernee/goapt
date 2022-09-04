package pkg

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type DebianSourceList []DebianSource

type DebianSourceType string

const (
	DebianSourceTypeDeb    = "deb"
	DebianSourceTypeDebSrc = "deb-src"
)

const (
	defaultSourceListPath = "/etc/apt/sources.list"
)

type DebianSource struct {
	Type      DebianSourceType
	URL       string
	Suite     string
	Component string
}

// DirectoryURL is the URL for the release metadata and directories
func (s *DebianSource) DirectoryURL() string {
	return fmt.Sprintf("%s/dists/%s/Release", s.URL, s.Suite)
}

// DirectorySignedURL is the URL for the release metadata and directories
func (s *DebianSource) DirectorySignedURL() string {
	return fmt.Sprintf("%s/dists/%s/InRelease", s.URL, s.Suite)
}

func (s *DebianSource) ResourceURL() string {
	var ret string

	switch s.Type {
	case DebianSourceTypeDeb:
		ret, _ = url.JoinPath(s.URL, "dists", s.Suite, s.Component, fmt.Sprintf("binary-%s", defaultArch), "Packages")
	case DebianSourceTypeDebSrc:
		ret, _ = url.JoinPath(s.URL, "dists", s.Suite, s.Component, "source", "Sources")
	}

	return ret
}

type Arch string

const (
	ArchAMD64 = "amd64"
	ArchI386  = "i386"
)

var defaultArch Arch = ArchAMD64

func SetArch(arch Arch) {
	defaultArch = arch
}

func LoadDebianSourceList() (DebianSourceList, error) {
	return loadDebianSourceList(defaultSourceListPath)
}

func loadDebianSourceList(path string) (DebianSourceList, error) {
	ret, err := loadDebianSourceFromFile(path)
	if err != nil {
		return ret, err
	}

	dir := path + ".d"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ret, nil
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		list, err := loadDebianSourceFromFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		ret = append(ret, list...)
	}
	return ret, nil
}

// loadDebianSourceFromFile loads one debian source.list file
//
// debian source format is:
//  deb [ option1=value1 option2=value2 ] uri suite [component1] [component2] [...]
//  deb-src [ option1=value1 option2=value2 ] uri suite [component1] [component2] [...]
// Example:
//  deb http://archive.ubuntu.com/ubuntu/ focal-updates main restricted
func loadDebianSourceFromFile(path string) (DebianSourceList, error) {
	var (
		ret DebianSourceList
		err error
	)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		lb, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		line := strings.TrimSpace(string(lb))
		if sList := parseLine(line); sList != nil {
			ret = append(ret, sList...)
		}
	}

	return ret, nil
}

func parseLine(line string) DebianSourceList {
	var (
		options []string
		args    []string
		ret     DebianSourceList
		splits  = strings.Split(line, " ")
	)

	if strings.HasPrefix(line, "#") {
		return ret
	}

	// Skip unknown types
	if len(splits) < 4 || splits[0] != "deb" && splits[0] != "deb-src" {
		return ret
	}

	// Split the options list and args list
	for _, split := range splits {
		if strings.Index(split, "=") != -1 {
			options = append(options, split)
		} else {
			args = append(args, split)
		}
	}

	if len(args) < 4 {
		return ret
	}

	// XXX: Options are ignored

	// Create a DebianSource struct for each component
	for _, component := range args[3:] {
		ret = append(ret, DebianSource{
			Type:      DebianSourceType(args[0]),
			URL:       args[1],
			Suite:     args[2],
			Component: component,
		})
	}

	return ret
}
