package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/GeertJohan/go.rice"
)

type Template struct {
	Filename string
	Contents []string
}

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

func main() {

	// Subcommands
	initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	tmplCommand := flag.NewFlagSet("template", flag.ExitOnError)

	// Init subcommand flag pointers
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	dirPtr := initCommand.String("o", dir, "Directory to store the mdd database")
	f := initCommand.Lookup("o")
	f.DefValue = fmt.Sprintf("The current path ie: '%s'", dir)

	base := filepath.Base(dir)
	projectPtr := initCommand.String("p", base, "Project name")
	f = initCommand.Lookup("p")
	f.DefValue = fmt.Sprintf("The current directory name ie: '%s'", base)

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		helptext := `
mdd is a tool for managing markdown system documentation

Usage:

	mdd <command> [arguments]

The commands are:

	init        initialise a mdd repository
	templates   list the templates available for use
`

		fmt.Println(helptext)
		os.Exit(1)
	}
	// Switch on the subcommand
	// Parse the flags for appropriate FlagSet
	// FlagSet.Parse() requires a set of arguments to parse as input
	// os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	switch os.Args[1] {
	case "init":
		initCommand.Parse(os.Args[2:])
	case "templates":
		tmplCommand.Parse(os.Args[2:])
	// case "new":
	// 	newCommand.Parse(os.Args[2:])
	// case "link":
	// 	linkCommand.Parse(os.Args[2:])
	// case "unlink":
	// 	unlinkCommand.Parse(os.Args[2:])
	// case "tags":
	// 	tagsCommand.Parse(os.Args[2:])
	// case "tag":
	// 	tagCommand.Parse(os.Args[2:])
	// case "untag":
	// 	untagCommand.Parse(os.Args[2:])
	// case "ls":
	// 	lsCommand.Parse(os.Args[2:])
	// case "verify":
	// 	verifyCommand.Parse(os.Args[2:])
	// case "server":
	// 	serverCommand.Parse(os.Args[2:])
	// case "publish":
	// 	publishCommand.Parse(os.Args[2:])
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Check which subcommand was Parsed using the FlagSet.Parsed() function. Handle each case accordingly.
	// FlagSet.Parse() will evaluate to false if no flags were parsed (i.e. the user did not provide any flags)
	if initCommand.Parsed() {
		// Asked for help
		if len(os.Args[2:]) > 0 && os.Args[2:][0] == "help" {
			helptext := `
mdd init creates a new mdd document repository

Usage:

	mdd init [arguments]

The arguments are:
`
			fmt.Println(helptext)
			initCommand.PrintDefaults()
			os.Exit(0)
		}
		_, err = NewProject(dirPtr, projectPtr)
	} else if tmplCommand.Parsed() {
		// Asked for help
		if len(os.Args[2:]) > 0 && os.Args[2:][0] == "help" {
			helptext := `
mdd templates lists the templates available

Usage:

	mdd templates
`
			fmt.Println(helptext)
			os.Exit(0)
		}
		p, err := FindProjectBelowCwd()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		for _, t := range p.Templates {
			log.Println("%s: %v", t.Filename, t.Contents)
		}

	}

	// Exit non zero on errrr
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func NewTemplate(path string) (Template, error) {
	t := Template{Filename: path}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return t, err
	}
	t.Contents = strings.Split(string(content), "\n")
	return t, nil
}

// FindProjectBelowCwd `walk`s the directory tree from '.' looking for the '.mdd' directory
func FindProjectBelowCwd() (Project, error) {

	p := Project{}

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == RootDirectory {
			p.HomePath = path
		}
		return nil
	})
	if err != nil {
		return p, err
	}

	p.TemplatePath = path.Join(p.HomePath, "templates")
	p.DataPath = path.Join(p.HomePath, "data")
	p.PublishPath = path.Join(p.HomePath, "publish")

	err = filepath.Walk(p.TemplatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Ignore failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return filepath.SkipDir
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
	if err := p.WriteProjectFile(map[string]string{"project": *name}); err != nil {
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
	return p, err
}

func (p *Project) ReadProjectFile() (map[string]string, error) {
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

func (p *Project) WriteProjectFile(db map[string]string) error {
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
