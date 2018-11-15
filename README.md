# mdd

`mdd` - Markdown documentation, a tool for capturing software engineering project documentation inside your git repository.

Why do you need `mdd`.

- Use `mdd` to capture your project requirements, architectural decisions, tests, and important meetings as [Markdown](https://en.wikipedia.org/wiki/Markdown) documentation, that can be versioned alongside your code in your source code control system.
- `mdd` comes with a set of templates to get you started.
- Link documents to one another, to make the relationships between then explicit. For example link a meetings to the requirements that come out of it, or link a requirements to its tests.
- Tag documents to help group related documents together eg: `#security` or `#performance`.
- Visualise the project documentation as a website, with hyperlinks to navigate and view the documentation as a whole.
- Markdown format works well with [Source code control systems](https://en.wikipedia.org/wiki/Source_Code_Control_System) like [git](https://en.wikipedia.org/wiki/Git). This means that:
  - Your documentation can be versioned alongside your code.
  - You can go back in time and see requirement changes.
  - No special tools are required to edit Markdown documents, thus making it suitable for use by all team members regarless of if they are developers, business analysts, testers and project managers.

# Installation

To build from source, you need a recent [go](http://golang.org) compiler installed, then run:

```
make build
```

To test, you need the [bats](https://github.com/sstephenson/bats) testing tool installed, then run:

```
make test
```

To install run:

```
sudo make install
```

**TODO** Add precompiled executables for OS X, Windows, Linux
 
# Architecture

`mdd` is a single executable file, that manages all functions of the system. The templates are compiled into the executable.

All the data is stored under a single directory `/.mdd`, and has the following structure. All files under .mdd should be added to your source code repository, with the **possible exception** of `./mdd/publish` because that directory contains the documents converted from Markdown to HTML format.

```
./.mdd
├── documents
├── project.data
├── publish
└── templates
    ├── adr.md
    ├── att.md
    ├── index.html
    ├── itst.md
    ├── mtg.md
    ├── nfr.md
    └── req.md
```

`mdd` stores metadata in two places. Project metadata is storted in a text file `./mdd/project.data`, and each Markdown documents contains its own metadata inside an HTML comment block.
When editing markdown documents, do not modify any of the lines that look like this, it
contains the metadata and are managed by the `mdd` command:

```
<!-- mdd
mdd-date-time: 2018-04-27
mdd-author: Fred
mdd-child: dec-AB0034
mdd-tag: front-end
mdd-tag: security
...
-->
```

# Usage

Use the command line tool `mdd`

`mdd help` displays basic help, and `mdd help command` displays help about a specific command.
 
The basic usage is as follows:

## Initialise

Create a new markdown document database, inside a git project `my-project`:

```
$ cd my-project
$ ls

$ mdd init my-project
$ ls -a
.mdd
```

## List templates

List the templates we have at our disposal:

```
$ mdd templates
   adr: Architecture Decision Record
   att: Automated test
  itst: Inspection test
   mtg: Meeting
   nfr: Non Functional Requirement
   req: Functional Requirement
```

# Create a document
 
Create a new document, specifying the template and document title:

_Note: the `-e` option opens the file in our `$EDITOR`_

```
$ mdd new req -e "User login"
mdd_docs/req/req-e2c-0001.md
```

The `e2c` is an (example) prefix that is generated from the users operating system username
to help eliminate a filename clash when multiple users are adding documents concurrently to the
same repository.
 
 
Add the `mdd_docs/*` files to your SCCS repository as you would for any other
text file.

Lets create a second document, for example a test script:

```
$ mdd new itst "Testing user login"
.mdd/documents/itst-b7-0002.md
```


To list documents with their title

```
$ mdd ls
itst-b7-0002.md       Testing user login
req-b7-0001.md        User login
```
 

Now we want to link the requirement to the test document to show they are related. We do that as follows: 

```
$ mdd link req-b7-0001 itst-b7-0002
req-b7-0001.md -> itst-b7-0002.md
```

Note: If we wanted to unlink the documents run the same command replacing `link` with `unlink`.

Tags are added and removed from documents using the mdd `tag` and `untag` commands eg:

 
```
$ mdd tag req-b7-0001 security
```

To display a summary of the links and tags against a documents, run the `ls -l` command, eg::

```
$ mdd ls -l
itst-b7-0002.md       Testing user login
req-b7-0001.md        User login                     #security
  -> itst-b7-0002.md  Testing user login

```

Verify the structure of the mdd database, use the `verify` command which will check:

-   Every links points to a valid document
-   Each `*[mdd-...]` section is valid syntactically

 
eg:

```
$ mdd verify
echo $?
0
```


To publish the database as html run:
 
```
$ mdd publish
$ open ./.mdd/publish/index.html
```

# Resources:

-   <https://stackoverflow.com/questions/44215896/markdown-metadata-format>
-   <https://github.com/jacebrowning/doorstop>
-   <https://github.com/rjeczalik/fs>
-   <https://github.com/russross/blackfriday>
- https://blog.rapid7.com/2016/08/04/build-a-simple-cli-tool-with-golang/
- http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
- https://github.com/sstephenson/bats
