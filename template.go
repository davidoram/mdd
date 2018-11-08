package main

import (
	"fmt"
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

// This structure is used for templates output
type TemplateView struct {
	Title string
}

var templateDesc *regexp.Regexp

func init() {
	templateDesc = regexp.MustCompile("^[# ]*([\\w-. ~]+) *$")

}

func (t *Template) ForView() TemplateView {
	return TemplateView{
		Title: t.Title,
	}
}

func ReadTemplate(path string) (Template, error) {
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

	if t.Title == "" {
		return t, fmt.Errorf("Template '%s', is missing a title", path)
	}

	return t, nil
}
