package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
)

type Document struct {
	Filename string
	Template *Template
	Title    string
}

var (
	titleRegex    *regexp.Regexp
	filenameRegex *regexp.Regexp
)

func init() {
	titleRegex = regexp.MustCompile("^[# ]*([\\w-. ~]+) *$")
	filenameRegex = regexp.MustCompile("^(\\w+)-(\\w+)-(\\d+)\\.md$")
}

func (d *Document) BaseFilename() string {
	return filepath.Base(d.Filename)
}

func (p *Project) ReadDocument(path string) (*Document, error) {

	d := Document{
		Filename: path,
	}

	base := filepath.Base(path)
	matches := filenameRegex.FindStringSubmatch(base)
	if len(matches) != 4 {
		return nil, fmt.Errorf("Filename '%s', doesnt match regex", base)
	}

	for _, t := range p.Templates {
		if t.Shortcut == matches[1] {
			d.Template = &t
			break
		}
	}
	if d.Template == nil {
		return nil, fmt.Errorf("No template for shortcode '%s'", matches[0])
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return &d, err
	}
	contents := strings.Split(string(content), "\n")

	// First line that matches the regex, is the description
	for _, l := range contents {
		matches := titleRegex.FindStringSubmatch(l)
		if len(matches) > 1 {
			d.Title = matches[1]
			break
		}
	}
	return &d, nil

}

func (p *Project) NewDocument(t *Template, title string) (Document, error) {

	d := Document{
		Filename: filepath.Join(p.DocumentPath, p.GenerateFilename(t)),
		Template: t,
		Title:    title,
	}

	f, err := os.Create(d.Filename)
	if err != nil {
		return d, err
	}
	defer f.Close()

	// Don't replace the title unless a new one supplied
	replacedTitle := false || title == ""

	for _, l := range t.Contents {
		if !replacedTitle {
			if titleRegex.MatchString(l) {
				l = fmt.Sprintf("# %s", d.Title)
				replacedTitle = true
			}
		}
		_, err = f.WriteString(fmt.Sprintf("%s\n", l))
		if err != nil {
			return d, err
		}
	}

	return d, nil
}

// GenerateFilename finds the next free filename for a given template
// and injects 3 chars from the USER envar to minimise classhes
func (p *Project) GenerateFilename(t *Template) string {
	max := 0
	err := filepath.Walk(p.DocumentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Ignore failure accessing a path %q: %v\n", path, err)
			return err
		}

		// Skip directories
		if info.IsDir() && path != p.DocumentPath {
			log.Print("skipping")
			return filepath.SkipDir
		}
		base := filepath.Base(path)
		matches := filenameRegex.FindStringSubmatch(base)
		if len(matches) == 4 {
			// filenames = append(filenames, base)
			var i int
			n, err := fmt.Sscanf(matches[3], "%d", &i)
			if n == 1 && err == nil {
				if i > max {
					max = i
				}
			}
		}
		return nil
	})
	// Turn username into a semi-unique part of the filename
	// to help avoid filename clashes
	u, err := user.Current()
	if err != nil || u.Username == "" {
		u = &user.User{Username: "unknown"}
	}
	h := md5.New()
	io.WriteString(h, u.Username)

	filename := fmt.Sprintf("%s-%x-%04d.md", t.Shortcut, h.Sum(nil)[0:1], max+1)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return filename
	}
	panic(fmt.Sprintf("GenerateFilename returning a file that already exists: '%s'", filename))
}
