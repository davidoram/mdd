PREFIX = /usr/local

.PHONY: go-deps
build: go-deps
	rice embed-go
	go build

.PHONY: run
run:
	go run main.go

.PHONY: go-deps
go-deps:
	go get github.com/GeertJohan/go.rice
	go get github.com/GeertJohan/go.rice/rice
	go get github.com/microcosm-cc/bluemonday
	go get gopkg.in/russross/blackfriday.v2

.PHONY: test
test:
	bats *.bats

.PHONY: test-deps
test-deps:
	open https://github.com/sstephenson/bats

.PHONY: install
install:
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp mdd $(DESTDIR)$(PREFIX)/bin/mdd

.PHONY: uninstall
uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/mdd

