CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -s src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-markdown-sqlite
	cp *.go src/github.com/whosonfirst/go-whosonfirst-markdown-sqlite/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-markdown"
	@GOPATH=$(GOPATH) go get -u "github.com/mattn/go-sqlite3"
	@GOPATH=$(GOPATH) go install "github.com/mattn/go-sqlite3"

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
	@GOPATH=$(shell pwd) go build -o bin/wof-markdown-sqlite cmd/wof-markdown-sqlite.go
