package main

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type Template struct {
	Filename string
	Shortcut string
	Contents []string
	Title    string
}

var templateDesc *regexp.Regexp

func init() {
	templateDesc = regexp.MustCompile("^[# ]*([\\w-. ~]+) *$")

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

	// First line that matches the regex, is the description
	for _, l := range t.Contents {
		matches := templateDesc.FindStringSubmatch(l)
		if len(matches) > 1 {
			t.Title = matches[1]
			break
		}
	}
	return t, nil
}
