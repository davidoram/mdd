package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/GeertJohan/go.rice"
)

// Project contains the details about an mdd project
//
// ./tmp
// └── .mdd							<- HomePath
//     ├── documents		<- DocumentPath : Documents live in here
//     ├── project.data <- A textual database containing project meadata
//     ├── publish			<- PublishPath: Publish the documents as an HTML website here
//     └── templates		<- TemplatePath: All the template files available to this project
type Project struct {
	HomePath     string
	TemplatePath string
	DocumentPath string
	PublishPath  string
	Templates    []Template
	Documents    []*Document
}

const (
	RootDirectory = ".mdd"
	ProjectDbFile = "project.data"
)

var (
	box       *rice.Box
	fileRegex regexp.Regexp
)

func init() {
	box = rice.MustFindBox("templates")
}

func NewProject(projectDir, name *string) (Project, error) {

	// HomePath is projectdir/.mdd
	p := Project{HomePath: path.Join(*projectDir, RootDirectory)}

	// Check that projectDir exists
	stat, err := os.Stat(*projectDir)
	if os.IsNotExist(err) {
		return p, fmt.Errorf("no such directory '%s', aborting", *projectDir)
	}
	if !stat.IsDir() {
		return p, fmt.Errorf("expect a directory not a file: '%s', aborting", *projectDir)
	}

	// Create the HomePath directory, error if it already exists
	_, err = os.Stat(p.HomePath)
	if !os.IsNotExist(err) {
		return p, fmt.Errorf("project directory '%s' already exists, aborting", p.HomePath)
	} else if os.IsNotExist(err) {
		log.Printf("Creating dir: '%s'\n", p.HomePath)
		if err := os.MkdirAll(p.HomePath, os.ModePerm); err != nil {
			log.Printf("Error creating dir: '%s', %v\n", p.HomePath, err)
			return p, err
		}
	} else if err != nil {
		return p, err
	}

	p.TemplatePath = path.Join(p.HomePath, "templates")
	p.DocumentPath = path.Join(p.HomePath, "documents")
	p.PublishPath = path.Join(p.HomePath, "publish")

	// Create our project database file
	log.Printf("Writing project database file")
	if err := p.writeProjectDb(map[string]string{"project": *name}); err != nil {
		log.Printf("Error writing project file: '%v'\n", err)
		return p, err
	}

	// Create directories for templates, data & publish
	log.Printf("Creating sub-directories")
	if err := os.MkdirAll(p.TemplatePath, os.ModePerm); err != nil {
		log.Printf("Error creating dir: '%s', %v\n", p.TemplatePath, err)
		return p, err
	}
	if err := os.MkdirAll(p.DocumentPath, os.ModePerm); err != nil {
		log.Printf("Error creating dir: '%s', %v\n", p.DocumentPath, err)
		return p, err
	}
	if err := os.MkdirAll(p.PublishPath, os.ModePerm); err != nil {
		log.Printf("Error creating dir: '%s', %v\n", p.PublishPath, err)
		return p, err
	}

	// Save a copy of all the templates
	log.Printf("Copying templates")
	err = box.Walk(".", func(pth string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("failure accessing a path inside box %q: %v\n", pth, err)
			return err
		}
		if !info.IsDir() {
			log.Printf("\t%s", pth)
			// read the whole file at once
			b, err := box.Bytes(pth)
			if err != nil {
				log.Printf("Error Reading template: '%s' from box, %v\n", pth, err)
				return err
			}

			// write the whole body at once
			tmplPath := path.Join(p.TemplatePath, pth)
			err = ioutil.WriteFile(tmplPath, b, 0644)
			if err != nil {
				log.Printf("Error writing template: '%s' to path '%s', %v\n", pth, tmplPath, err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return p, err
	}
	// Return a correctly initialised Project structure
	return ReadProject(p.HomePath, false)
}

// Context provided to the Project directory walk fn
type projectWalkCtx struct {
	HomePaths []string
}

var errNoProjectFound = fmt.Errorf("No project found")

// FindProjectBelowCwd `walk`s the directory tree from '.' looking for the '.mdd' directory
// If it finds a project, will return it, otherwise will return errNoProjectFound
// If ignoreBrokenFiles is true will skip over broken Templates and Documents, which is
// useful for being able to work on projects that have small problems
func FindProjectBelowCwd(ignoreBrokenFiles bool) (*Project, error) {
	ctx := projectWalkCtx{HomePaths: make([]string, 0)}
	//log.Printf("Looking for project ...")
	err := filepath.Walk(".", ctx.projectWalkFn)
	if err != nil {
		return nil, err
	}

	if len(ctx.HomePaths) > 0 {
		// Return a correctly initialised Project structure
		p, err := ReadProject(ctx.HomePaths[0], ignoreBrokenFiles)
		if err != nil {
			return nil, err
		}
		return &p, err
	}
	return nil, errNoProjectFound
}

func (ctx *projectWalkCtx) projectWalkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		//log.Printf("error accessing a path %q: %v\n", path, err)
		return err
	}
	if info.IsDir() && info.Name() == RootDirectory {
		ctx.HomePaths = append(ctx.HomePaths, path)
		//log.Printf("Found project at '%s'", path)
		return filepath.SkipDir
	}
	return nil
}

func ReadProject(homePath string, ignoreBrokenFiles bool) (Project, error) {

	p := Project{HomePath: homePath}
	p.TemplatePath = path.Join(p.HomePath, "templates")
	p.DocumentPath = path.Join(p.HomePath, "documents")
	p.PublishPath = path.Join(p.HomePath, "publish")

	// Read the templates
	err := filepath.Walk(p.TemplatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Ignore failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Only treat markdown files as templates
		if filepath.Ext(path) == ".md" {
			tmpl, err := ReadTemplate(path)
			if err != nil {
				if !ignoreBrokenFiles {
					return err
				}
				log.Printf("Ignoring template at path %q: %v\n", path, err)
			} else {
				p.Templates = append(p.Templates, tmpl)
			}
		}
		return nil
	})
	if err != nil {
		return p, err
	}

	// Read the documents
	err = filepath.Walk(p.DocumentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Ignore failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		doc, err := p.ReadDocument(path)
		if err != nil {
			if !ignoreBrokenFiles {
				return err
			}
			log.Printf("Ignoring document at path %q: %v\n", path, err)
		} else {
			p.Documents = append(p.Documents, doc)
		}
		return nil
	})
	if err != nil {
		return p, err
	}

	return p, err
}

func (p *Project) readProjectDb() (map[string]string, error) {
	db := map[string]string{}
	dbPath := path.Join(p.HomePath, ProjectDbFile)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return db, nil
	}

	f, err := os.Open(dbPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// The regex to match
	// key: value one
	// in our file
	re := regexp.MustCompile("^[[:space:]]*([[:word:]]-)+:[[:space:]]+([[:word:]] -)+[[:space:]]*$")
	r := bufio.NewReader(f)
	s, e := r.ReadString('\n')
	for e == nil {
		matches := re.FindAllString(s, -1)
		if len(matches) == 0 {
			db[matches[0]] = matches[1]
		}
		s, e = r.ReadString('\n')
	}
	return db, nil
}

func (p *Project) writeProjectDb(db map[string]string) error {
	dbPath := path.Join(p.HomePath, ProjectDbFile)
	f, err := os.Create(dbPath)
	if err != nil {
		return err
	}
	f.WriteString("# mdd project db file. Do not edit\n")
	for k, v := range db {
		_, err := f.WriteString(fmt.Sprintf("%s: %s\n", k, v))
		if err != nil {
			return err
		}
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) DeleteAllPublished() error {
	// Open the directory and read all its files.
	dirRead, err := os.Open(p.PublishPath)
	if err != nil {
		return err
	}
	dirFiles, err := dirRead.Readdir(0)
	if err != nil {
		return err
	}

	// Loop over the directory's files.
	for _, f := range dirFiles {

		// Get name of file and its full path.
		fullPath := filepath.Join(p.PublishPath, f.Name())

		// Remove the file.
		err = os.Remove(fullPath)
		if err != nil {
			return err
		}
	}
	return nil
}
