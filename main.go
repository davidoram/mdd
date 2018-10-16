package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func init() {
	// Remove the Date/Time from log messages
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
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
		if len(p.Templates) == 0 {
			log.Printf("Project %s, has no templates", p.HomePath)
		}
		for _, t := range p.Templates {
			log.Printf("%6s: %s", t.Shortcut, t.Description)
		}

	}

	// Exit non zero on errrr
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
