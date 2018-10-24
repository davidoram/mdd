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

const (
	MetadataSeparator = ":"
	MetadataChild     = "mdd-child"
	MetadataTag       = "mdd-tag"
	MetadataStart     = "<!-- mdd"
	MetadataEnd       = "-->"
)

type Document struct {
	Filename string
	Template *Template
	Title    string

	// Metadata
	Children map[string]bool
	Tags     map[string]bool

	// File contents
	raw []byte
}

var (
	titleRegex     *regexp.Regexp
	filenameRegex  *regexp.Regexp
	metaStartRegex *regexp.Regexp
	metaEndRegex   *regexp.Regexp
	tagRegex       *regexp.Regexp
)

func init() {
	titleRegex = regexp.MustCompile("^[# ]*([\\w-. ~]+) *$")
	filenameRegex = regexp.MustCompile("^(\\w+)-(\\w+)-(\\d+)\\.md$")
	metaStartRegex = regexp.MustCompile("^\\s*<!-- mdd\\s*$")
	metaEndRegex = regexp.MustCompile("^\\s*-->\\s*$")
	tagRegex = regexp.MustCompile("^[[:word:]]{3,20}$")
}

func (d *Document) BaseFilename() string {
	return filepath.Base(d.Filename)
}

func (d *Document) AddChild(child *Document) error {
	d.Children[child.BaseFilename()] = true
	return nil
}

func (d *Document) RemoveChild(childFilename string) error {
	delete(d.Children, childFilename)
	return nil
}

func (d *Document) Tag(tag string) error {
	if !tagRegex.MatchString(tag) {
		return fmt.Errorf("Tags must be 3-20 chars long, made up of the following characters: '0-9A-Za-z_'")
	}
	d.Tags[tag] = true
	return nil
}

func (d *Document) Untag(tag string) error {
	delete(d.Tags, tag)
	return nil
}

func (d *Document) TagNames() []string {
	tags := []string{}
	for tag, _ := range d.Tags {
		tags = append(tags, tag)
	}
	return tags
}

func (p *Project) ReadDocument(path string) (*Document, error) {

	d := Document{
		Filename: path,
		Children: make(map[string]bool),
		Tags:     make(map[string]bool),
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

	var err error
	d.raw, err = ioutil.ReadFile(path)
	if err != nil {
		return &d, err
	}
	content := string(d.raw)
	contents := strings.Split(content, LineBreak)

	// First line that matches the regex, is the description
	for _, l := range contents {
		matches := titleRegex.FindStringSubmatch(l)
		if len(matches) > 1 {
			d.Title = matches[1]
			break
		}
	}

	// First line that matches the regex, is the description
	inMeta := false
	for _, l := range contents {
		if inMeta {
			if metaEndRegex.MatchString(l) {
				// log.Printf("%d end meta", i)
				inMeta = false
			} else {
				// log.Printf("%d in meta", i)
				d.parseMetadata(l)
			}
		} else {
			if metaStartRegex.MatchString(l) {
				// log.Printf("%d start meta", i)
				inMeta = true
			}
		}
	}
	return &d, nil
}

func (d *Document) WriteDocument() error {
	file, err := os.OpenFile(d.Filename, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	s := string(d.raw)

	// Construct the new metadata
	meta := strings.Join(d.metadataForWrite(), LineBreak)

	// Replace the old metadata
	r := regexp.MustCompile("(?ms)^<!-- mdd$(.*)^-->$")
	s = r.ReplaceAllString(s, meta)

	// Write the file
	_, err = file.WriteString(s)
	if err != nil {
		return err
	}

	return nil
}

// Return the metadata as an array suitable for writing out to the file
func (d *Document) metadataForWrite() []string {
	meta := []string{MetadataStart}

	for key := range d.Children {
		meta = append(meta, fmt.Sprintf("%s: %s", MetadataChild, key))
	}
	for key := range d.Tags {
		meta = append(meta, fmt.Sprintf("%s: %s", MetadataTag, key))
	}
	meta = append(meta, MetadataEnd)
	return meta
}

// line has one of the forms:
// mdd-child:document-name
// mdd-tag:value
func (d *Document) parseMetadata(line string) error {

	meta := strings.Split(line, MetadataSeparator)
	if len(meta) != 2 {
		return fmt.Errorf("Expected 2 values, found %d from metadata '%s'", len(meta), line)
	}

	key := strings.TrimSpace(meta[0])
	value := strings.TrimSpace(meta[1])
	switch key {
	case MetadataChild:
		d.Children[value] = true
	case MetadataTag:
		d.Tags[value] = true
	default:
		return fmt.Errorf("Unrecognised metadata tag '%s'", key)
	}
	return nil
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
