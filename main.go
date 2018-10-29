package main

import (
	"flag"
	"fmt"
	"html/template"
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
	ls					list documents created
	link				link a parent and child document
	unlink			remove the link between a parent and child document
	tag					tag a document
	untag				untag a document
	verify			verify the struture of the mdd repository documents
	publish			create a static website reflectings the mdd repository
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
	lsCommand := flag.NewFlagSet("ls", flag.ExitOnError)
	linkCommand := flag.NewFlagSet("link", flag.ExitOnError)
	unlinkCommand := flag.NewFlagSet("unlink", flag.ExitOnError)
	tagCommand := flag.NewFlagSet("tag", flag.ExitOnError)
	untagCommand := flag.NewFlagSet("untag", flag.ExitOnError)
	verifyCommand := flag.NewFlagSet("verify", flag.ExitOnError)
	publishCommand := flag.NewFlagSet("publish", flag.ExitOnError)

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

	linkPtr := lsCommand.Bool("l", false, "Display links")

	publishPtr := publishCommand.String("o", dir, "Directory to publish the site to, defaults .mdd/publish")

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
	case "ls":
		lsCommand.Parse(os.Args[2:])
		err = doLs(lsCommand, linkPtr)
	case "link":
		if len(os.Args) >= 3 {
			linkCommand.Parse(os.Args[2:])
			err = doLink(linkCommand)
		} else {
			err = fmt.Errorf("Cannot parse command line. Try 'mdd link help'")
		}
	case "unlink":
		if len(os.Args) >= 3 {
			unlinkCommand.Parse(os.Args[2:])
			err = doUnlink(unlinkCommand)
		} else {
			err = fmt.Errorf("Cannot parse command line. Try 'mdd unlink help'")
		}
	case "tag":
		if len(os.Args) >= 3 {
			tagCommand.Parse(os.Args[2:])
			err = doTag(tagCommand)
		} else {
			err = fmt.Errorf("Cannot parse command line. Try 'mdd tag help'")
		}
	case "untag":
		if len(os.Args) >= 3 {
			untagCommand.Parse(os.Args[2:])
			err = doUntag(untagCommand)
		} else {
			err = fmt.Errorf("Cannot parse command line. Try 'mdd untag help'")
		}
	case "verify":
		verifyCommand.Parse(os.Args[2:])
		err = doVerify(verifyCommand)

	case "publish":
		publishCommand.Parse(os.Args[2:])
		err = doPublish(publishCommand, publishPtr)

	// case "server":
	// 	serverCommand.Parse(os.Args[2:])
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

func doLs(flags *flag.FlagSet, linkPtr *bool) error {
	helptext := `
mdd ls lists all the documents created

Usage:

	mdd ls

The arguments are:
`
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

	p, err := FindProjectBelowCwd()
	if err != nil {
		return err
	}
	for _, d := range p.Documents {
		if len(d.TagNames()) > 0 {
			log.Printf("%-15s  %-30s %s", d.BaseFilename(), d.Title, d.TagNames())
		} else {
			log.Printf("%-15s  %-30s", d.BaseFilename(), d.Title)
		}

		// Display links?
		if *linkPtr {
			for name := range d.Children {
				log.Printf("  -> %-15s", name)
			}
		}
	}
	return nil
}

func doLink(flags *flag.FlagSet) error {
	helptext := `
mdd link links a parent and child document

Usage:

	mdd link parent child

parent is the parent documents filename.
child is the child documents filename.

The arguments are:
`
	// log.Printf("%v %d", os.Args, len(os.Args))
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
	if len(os.Args[2:]) != 2 {
		fmt.Println(helptext)
		flags.PrintDefaults()
		return fmt.Errorf("Missing arguments")
	}

	parent := os.Args[2:][0]
	if !strings.HasSuffix(parent, ".md") {
		parent = fmt.Sprintf("%s.md", parent)
	}
	child := os.Args[2:][1]
	if !strings.HasSuffix(child, ".md") {
		child = fmt.Sprintf("%s.md", child)
	}
	if parent == child {
		return fmt.Errorf("Cant link to self")
	}

	p, err := FindProjectBelowCwd()
	if err != nil {
		return err
	}

	var pdoc *Document
	var cdoc *Document
	for _, d := range p.Documents {
		if pdoc == nil && d.BaseFilename() == parent {
			pdoc = d
		}
		if cdoc == nil && d.BaseFilename() == child {
			cdoc = d
		}
	}
	if pdoc == nil {
		return fmt.Errorf("Cant find parent '%s'", parent)
	}
	if cdoc == nil {
		return fmt.Errorf("Cant find child '%s'", child)
	}
	// log.Printf("OK %s -> %s", pdoc.BaseFilename(), cdoc.BaseFilename())
	if err = pdoc.AddChild(cdoc); err == nil {
		err = pdoc.WriteDocument()
	}
	return nil

}

func doUnlink(flags *flag.FlagSet) error {
	helptext := `
mdd unlink breaks the link between a parent and child document

Usage:

	mdd unlink parent child

parent is the parent documents filename.
child is the child documents filename.

The arguments are:
`
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

	// Missing document
	if len(os.Args[2:]) != 2 {
		fmt.Println(helptext)
		flags.PrintDefaults()
		return fmt.Errorf("Missing arguments")
	}

	parent := os.Args[2:][0]
	if !strings.HasSuffix(parent, ".md") {
		parent = fmt.Sprintf("%s.md", parent)
	}
	child := os.Args[2:][1]
	if !strings.HasSuffix(child, ".md") {
		child = fmt.Sprintf("%s.md", child)
	}
	if parent == child {
		return fmt.Errorf("Cant unlink from self")
	}

	p, err := FindProjectBelowCwd()
	if err != nil {
		return err
	}

	var pdoc *Document
	for _, d := range p.Documents {
		if pdoc == nil && d.BaseFilename() == parent {
			pdoc = d
		}
	}
	if pdoc == nil {
		return fmt.Errorf("Cant find parent '%s'", parent)
	}
	if err = pdoc.RemoveChild(child); err == nil {
		err = pdoc.WriteDocument()
	}
	return nil

}

func doTag(flags *flag.FlagSet) error {
	helptext := `
mdd tag adds tags to a document

Usage:

	mdd tag document tag tag2 ...

document is a documents filename.

The arguments are:
`
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

	// Missing document & tag
	if len(os.Args[2:]) < 2 {
		fmt.Println(helptext)
		flags.PrintDefaults()
		return fmt.Errorf("Missing arguments")
	}

	document := os.Args[2:][0]
	if !strings.HasSuffix(document, ".md") {
		document = fmt.Sprintf("%s.md", document)
	}

	tags := os.Args[3:]

	p, err := FindProjectBelowCwd()
	if err != nil {
		return err
	}

	var doc *Document
	for _, d := range p.Documents {
		if doc == nil && d.BaseFilename() == document {
			doc = d
		}
	}
	if doc == nil {
		return fmt.Errorf("Cant find parent '%s'", document)
	}
	for _, t := range tags {
		if err = doc.Tag(t); err != nil {
			return err
		}
	}
	return doc.WriteDocument()
}

func doUntag(flags *flag.FlagSet) error {
	helptext := `
mdd untag removes tags from a document

Usage:

	mdd untag document tag tag2 ...

document is a documents filename.

The arguments are:
`
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

	// Missing document
	if len(os.Args[2:]) != 2 {
		fmt.Println(helptext)
		flags.PrintDefaults()
		return fmt.Errorf("Missing arguments")
	}

	document := os.Args[2:][0]
	if !strings.HasSuffix(document, ".md") {
		document = fmt.Sprintf("%s.md", document)
	}

	tags := os.Args[3:]

	p, err := FindProjectBelowCwd()
	if err != nil {
		return err
	}

	var doc *Document
	for _, d := range p.Documents {
		if doc == nil && d.BaseFilename() == document {
			doc = d
		}
	}
	if doc == nil {
		return fmt.Errorf("Cant find parent '%s'", document)
	}
	for _, t := range tags {
		if err = doc.Untag(t); err != nil {
			return err
		}
	}
	return doc.WriteDocument()
}

func doVerify(flags *flag.FlagSet) error {
	helptext := `
mdd verify checks the integrity of the documention.

Usage:

	mdd verify

The return code will be zero if no errors exist, non-zero if one or more errors are detected.
mdd verify is suitable for injecting into a CI pipeline to verify that documentation meets the
basic level of structural checks.

The arguments are:
`
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

	p, err := FindProjectBelowCwd()
	if err != nil {
		return err
	}
	errors := 0
	for _, d := range p.Documents {
		// Check each child pointer is valid
		for name := range d.Children {
			found := false
			for _, c := range p.Documents {
				if c.BaseFilename() == name {
					found = true
					break
				}
			}
			if !found {
				errors += 1
				fmt.Printf("'%s' has child '%s' which doesnt exist\n", d.BaseFilename(), name)
			}
		}
	}

	if errors != 0 {
		return fmt.Errorf("Found %d errors", errors)
	}
	return nil
}

func doPublish(flags *flag.FlagSet, dirPtr *string) error {
	helptext := `
mdd publish creates a static website for the mdd repository

Usage:

	mdd publish [arguments]

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
	p, err := FindProjectBelowCwd()
	if err != nil {
		return err
	}

	err = p.DeleteAllPublished()
	if err != nil {
		return err
	}

	// Build up a data structure to use when we spit out the
	// index.html file
	data := struct {
		// Map from tags -> DocView with that tag
		TagDocs map[string][]DocView

		// Map from Template filename to Documents following that template
		TmplDocs map[string][]DocView

		// Map from Template filename to Template title
		TmplTitles map[string]string
	}{
		TagDocs:    make(map[string][]DocView),
		TmplDocs:   make(map[string][]DocView),
		TmplTitles: make(map[string]string),
	}

	for _, d := range p.Documents {

		dv := d.ForView()

		// Index by Tag
		for _, t := range d.TagNames() {
			if data.TagDocs[t] == nil {
				data.TagDocs[t] = make([]DocView, 0)
			}
			data.TagDocs[t] = append(data.TagDocs[t], dv)
		}

		// Index by Template
		if data.TmplDocs[dv.TemplateFilename] == nil {
			data.TmplDocs[dv.TemplateFilename] = make([]DocView, 0)
			data.TmplTitles[dv.TemplateFilename] = dv.TemplateTitle
		}
		data.TmplDocs[dv.TemplateFilename] = append(data.TmplDocs[dv.TemplateFilename], d.ForView())

		log.Printf("Converting %s\n", d.Filename)
		err = d.ConvertToHTML(p.PublishPath)
		if err != nil {
			return err
		}

	}

	// Create the index.html document
	tmpl, err := template.ParseFiles(filepath.Join(p.TemplatePath, "index.html"))
	if err != nil {
		return err
	}

	outFile := filepath.Join(p.PublishPath, "index.html")
	_, err = os.Stat(outFile)
	if !os.IsNotExist(err) {
		return fmt.Errorf("index file '%s' already exists", outFile)
	}
	file, err := os.OpenFile(outFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}
	log.Printf("%s\n", outFile)
	return nil
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
