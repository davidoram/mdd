# mdd
Markdown docs


 

-   Templated documents markdown + metadata
-   Metadata classifies & links documents
-   Templates provide starting structure to capture different documentation

    -   Architectural decisions
    -   Meeting minutes
    -   Requirements
    -   Tests
    -   ... your own
-   Links provide a way to show the linkages
    -   Requirement X is verified with Test y
    -   maintained through the mdd tool
-   Verification
    -   Checks that all links exist
    -   Incorporate into CI/CD
        -   All links valid
        -   Others eg: All Requirements covered by tests

-   Documentation
    -   Run as website
    -   Spit out a website
    -   Summary pages
        -   TOC

 

Resources:

-   <https://stackoverflow.com/questions/44215896/markdown-metadata-format>
-   <https://github.com/jacebrowning/doorstop>
-   <https://github.com/rjeczalik/fs>
-   <https://github.com/russross/blackfriday>
- https://blog.rapid7.com/2016/08/04/build-a-simple-cli-tool-with-golang/

 

 

Interface
=========

Use the command line tool `mdd`

 

Create a new markdown document database, inside a github project:


```
$ cd my-project
$ ls

$ mdd init my-project
$ ls -a
.mdd

$ cat .mdd
mdd-project: my-project
...
```

This will create the `mdd_docs` directory

List the templates we have at our disposal:

```$ mdd tmpl
req     Requirements
tst     Test
mtg     Meeting minutes
dec     Decision
....
```
 
Create a new document from a template:



```$ mdd new req "User login"
mdd_docs/req/req-AAA0001
# Will open the file if $EDITOR set
```

The `AAA` is a  prefix that is generated from the `${USER}` environment variable
to minimize clashes with multiple users adding records concurrently

 

When editing the file, do not modify any of the lines that look like this. They
contain the metadata and are managed by the `mdd` command.

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

 
Add the `mdd_docs/*` files to your git repo as you would for any other
documentation.

 

To link two documents, for example if we have a decision `dec-GHF0035` relating
to a requirement `req-AA0001` then we say that the decision is a child of the
requirement as follows:

 

```
$ mdd link req-AA0001 dec-GHF0035
  req-AA0001   User login
+ dec-GHF0035  Support OAUTH2

```

 

To unlink two documents, use the following:

 

```
$ mdd unlink req-AA0001 dec-GHF0035
  req-AA0001   User login
- dec-GHF0035  Support OAUTH2

```

 

Tags are added and removed as follows

 

```
$ mdd tags
tags:
  - front-end
  - back-end
  - performance
  - security

$ mdd tag dec-GHF0035 front-end
dec-GHF0035       Support OAUTH2
tags:
  - back-end
  - front-end
  - security

$ mdd untag dec-GHF0035 back-end
dec-GHF0035       Support OAUTH2
tags:
  - front-end
  - security

```

 

To display a summary of the links and tags against a documents, use the
following:

 

```
$ mdd ls dec-GHF0035
dec-GHF0035       Support OAUTH2
tags:
  - front-end
  - security
parents:
  - req-AAA0001   User Login
childeren:
  - tst-SFS0056   OAUTH tests
```

 

To verify the structure of the mdd database. It checks:

-   Every links points to a valid document
-   Each `*[mdd-...]` section is valid syntactically

 

```
$ mdd verify
# exit code 0 means ok, !=0 means error
```

 

To serve up the documentation

 

```
$ mdd server [-p port]
Listening on 6061

open http://localhost:6061
```


To publish the documentation as a static website

 

```
$ mdd publish ./doc-site
$ ls doc-site/
index.html
...

```
 
<!-- mdd
*[mdd-tag]: security
-->
 
