package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	helptext = `
mdd is a tool for managing markdown system documentation

Usage:

	mdd <command> [arguments]

The commands are:

	init        initialise a mdd repository
	templates   list the templates available for use
	new         add a new document based on a template
`
)

func init() {
	// Remove the Date/Time from log messages
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

func main() {

	// Subcommands
	initCommand := flag.NewFlagSet("init", flag.ExitOnError)
	tmplCommand := flag.NewFlagSet("template", flag.ExitOnError)
	newCommand := flag.NewFlagSet("new", flag.ExitOnError)

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

	editPtr := newCommand.Bool("e", false, "Open the new file in your $EDITOR")

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
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
		err = doInit(initCommand, dirPtr, projectPtr)
	case "templates":
		tmplCommand.Parse(os.Args[2:])
		err = doTemplates(tmplCommand)
	case "new":
		if len(os.Args) >= 3 {
			newCommand.Parse(os.Args[3:])
			err = doNew(newCommand, editPtr)
		} else {
			err = fmt.Errorf("Cannot parse command line. Try 'mdd new help'")
		}

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
		log.Printf("Unknown command '%s'", os.Args[1])
		fmt.Println(helptext)
		os.Exit(1)
	}

	// Exit non zero on error
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

}

func doInit(flags *flag.FlagSet, dirPtr, projectPtr *string) error {
	helptext := `
mdd init creates a new mdd document repository

Usage:

	mdd init [arguments]

The arguments are:
`
	// FlagSet.Parse() will evaluate to false if no flags were parsed
	if !flags.Parsed() {
		return fmt.Errorf("Error parsing arguments")
	}

	// Asked for help?
	if len(os.Args[2:]) > 0 && os.Args[2:][0] == "help" {
		log.Println(helptext)
		flags.PrintDefaults()
		return nil
	}
	_, err := NewProject(dirPtr, projectPtr)
	return err
}

func doTemplates(flags *flag.FlagSet) error {
	helptext := `
mdd templates lists the templates available

Usage:

	mdd templates
`
	// FlagSet.Parse() will evaluate to false if no flags were parsed
	if !flags.Parsed() {
		return fmt.Errorf("Error parsing arguments")
	}

	// Asked for help
	if len(os.Args[2:]) > 0 && os.Args[2:][0] == "help" {
		fmt.Println(helptext)
		flags.PrintDefaults()
		return nil
	}
	p, err := FindProjectBelowCwd()
	if err != nil {
		return err
	}
	if len(p.Templates) == 0 {
		log.Printf("Project %s, has no templates", p.HomePath)
	}
	for _, t := range p.Templates {
		log.Printf("%6s: %s", t.Shortcut, t.Title)
	}
	return nil
}

func doNew(flags *flag.FlagSet, openEditor *bool) error {
	helptext := `
mdd new creates a new document from a template

Usage:

	mdd new template [arguments] [title]

template is the template to use for the new document.
title is an optional title for the document.

The arguments are:
`
	log.Printf("flag.,Args: %+v %d", flags.Args(), len(flags.Args()))
	// FlagSet.Parse() will evaluate to false if no flags were parsed
	if !flags.Parsed() {
		return fmt.Errorf("Error parsing arguments")
	}
	// Asked for help?
	if len(os.Args[2:]) > 0 && os.Args[2:][0] == "help" {
		fmt.Println(helptext)
		flags.PrintDefaults()
		return nil
	}

	// Missing template shortcut
	if len(os.Args[2:]) == 0 {
		return fmt.Errorf("Missing 'template shortcut' argument")
	}
	shortcut := os.Args[2:][0]
	p, err := FindProjectBelowCwd()
	if err != nil {
		return err
	}
	title := ""
	if len(flags.Args()) > 0 {
		title = flags.Args()[0]
	}
	for _, t := range p.Templates {
		if shortcut == t.Shortcut {
			doc, err := p.NewDocument(&t, title)
			if err != nil {
				return err
			}
			log.Printf("%s", doc.Filename)
			if *openEditor {
				return execEditor(doc.Filename)
			}
			return nil
		}
	}
	return fmt.Errorf("No such template: '%s'", shortcut)
}

func execEditor(filename string) error {
	val := ""
	val, ok := os.LookupEnv("EDITOR")
	if !ok {
		return fmt.Errorf("Envar EDITOR not set")
	}

	// EDITOR miight be set to a value like '/path/to/editor --some-flags', so we
	// need to parse the binary from the args
	editorArgs := strings.Split(val, " ")
	binary, err := exec.LookPath(editorArgs[0])
	if err != nil {
		return err
	}

	args := []string{}
	args = append(args, editorArgs[1:]...)
	args = append(args, filename)

	err = syscall.Exec(binary, args, os.Environ())
	return err
}
