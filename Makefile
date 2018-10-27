build: get-deps
	rice embed-go
	go build

run:
	go run main.go

get-deps:
	go get github.com/GeertJohan/go.rice
	go get github.com/GeertJohan/go.rice/rice
	go get github.com/microcosm-cc/bluemonday
	go get gopkg.in/russross/blackfriday.v2
