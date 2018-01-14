CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -s src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-markdown-search
	cp *.go src/github.com/whosonfirst/go-whosonfirst-markdown-search/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-markdown"
	@GOPATH=$(GOPATH) go get -u "github.com/blevesearch/bleve"
	@GOPATH=$(GOPATH) go get -u "github.com/mattn/go-sqlite3"
	@GOPATH=$(GOPATH) go install "github.com/mattn/go-sqlite3"
	if test ! -d src/gopkg.in; then mkdir -p src/gopkg.in; fi
	mv src/github.com/whosonfirst/go-whosonfirst-markdown/vendor/gopkg.in/russross src/gopkg.in
	mv src/github.com/whosonfirst/go-whosonfirst-markdown/vendor/github.com/shurcooL src/github.com/

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt *.go

bin: 	rmdeps self
	rm -rf bin/*
	@GOPATH=$(shell pwd) go build -o bin/wof-markdown-search cmd/wof-markdown-search.go
