package main

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type Template struct {
	Filename    string
	Shortcut    string
	Contents    []string
	Description string
}

var templateDescriptionRegex *regexp.Regexp

func init() {
	templateDescriptionRegex = regexp.MustCompile("^#(.+)$")

}

func NewTemplate(path string) (Template, error) {
	t := Template{Filename: path}

	// Shortcut is the Base minus extension
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	t.Shortcut = base[0 : len(base)-len(ext)]

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return t, err
	}
	t.Contents = strings.Split(string(content), "\n")

	// First non-empty line is the description
	re := regexp.MustCompile("^[# ]*([\\w-. ~]+) *$")

	for _, l := range t.Contents {
		matches := re.FindStringSubmatch(l)
		if len(matches) > 1 {
			t.Description = matches[1]
			break
		}
	}
	return t, nil
}
