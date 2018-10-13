package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

type Project struct {
	HomePath string
}

func main() {

	// Subcommands
	initCommand := flag.NewFlagSet("init", flag.ExitOnError)

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
	// case "tmpl":
	// 	templateCommand.Parse(os.Args[2:])
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
		p := Project{HomePath: *dirPtr}
		err = p.init(projectPtr)
	}
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func (p *Project) ReadProjectFile() (map[string]string, error) {
	db := map[string]string{}
	dbPath := path.Join(p.HomePath, ".mdd")
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
	dbPath := path.Join(p.HomePath, ".mdd")
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

func (p *Project) init(name *string) error {

	// Create the HomePath directory if it doesnt exist
	if _, err := os.Stat(p.HomePath); os.IsNotExist(err) {
		if err := os.MkdirAll(p.HomePath, os.ModePerm); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Write the project file, error if it already exists
	dbFile := path.Join(p.HomePath, ".mdd")
	_, err := os.Stat(dbFile)
	if !os.IsNotExist(err) {
		return fmt.Errorf("Wont overwrite existing project at '%s'", dbFile)
	}
	return p.WriteProjectFile(map[string]string{"project": *name})

}
