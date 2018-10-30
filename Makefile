build: go-deps
	rice embed-go
	go build

run:
	go run main.go

go-deps:
	go get github.com/GeertJohan/go.rice
	go get github.com/GeertJohan/go.rice/rice
	go get github.com/microcosm-cc/bluemonday
	go get gopkg.in/russross/blackfriday.v2


test:
	bats mdd.bats

test-deps:
	open https://github.com/sstephenson/bats
