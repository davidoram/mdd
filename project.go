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

type Project struct {
	HomePath     string
	TemplatePath string
	DataPath     string
	PublishPath  string
	Templates    []Template
}

const (
	RootDirectory = ".mdd"
	ProjectDbFile = "project.data"
)

var box *rice.Box

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
		return p, fmt.Errorf("project directory '%s' alredy exists, aborting", p.HomePath)
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
	p.DataPath = path.Join(p.HomePath, "data")
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
	if err := os.MkdirAll(p.DataPath, os.ModePerm); err != nil {
		log.Printf("Error creating dir: '%s', %v\n", p.DataPath, err)
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
	return ReadProject(p.HomePath)
}

// Context provided to the Project directory walk fn
type projectWalkCtx struct {
	HomePath string
}

// FindProjectBelowCwd `walk`s the directory tree from '.' looking for the '.mdd' directory
// If it finds a project, will return it, otherwise will return nil
func FindProjectBelowCwd() (*Project, error) {

	ctx := projectWalkCtx{}
	//log.Printf("Looking for project ...")
	err := filepath.Walk(".", ctx.projectWalkFn)
	if err != nil && ctx.HomePath != "" {
		return nil, err
	}
	if ctx.HomePath != "" {
		// Return a correctly initialised Project structure
		p, err := ReadProject(ctx.HomePath)
		if err != nil {
			return nil, err
		}
		return &p, err
	}
	return nil, err
}

func (ctx *projectWalkCtx) projectWalkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Printf("error accessing a path %q: %v\n", path, err)
		return err
	}
	if info.IsDir() && info.Name() == RootDirectory {
		ctx.HomePath = path
		// log.Printf("Found project at '%s'", path)
		return filepath.SkipDir
	}
	return nil
}

func ReadProject(homePath string) (Project, error) {

	p := Project{HomePath: homePath}
	p.TemplatePath = path.Join(p.HomePath, "templates")
	p.DataPath = path.Join(p.HomePath, "data")
	p.PublishPath = path.Join(p.HomePath, "publish")

	err := filepath.Walk(p.TemplatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Ignore failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		tmpl, err := NewTemplate(path)
		if err != nil {
			log.Printf("Ignoring template at path %q: %v\n", path, err)
		} else {
			p.Templates = append(p.Templates, tmpl)
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
