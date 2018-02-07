CWD=$(shell pwd)
GOPATH := $(CWD)

prep:
	if test -d pkg; then rm -rf pkg; fi

self:   prep
	if test -s src; then rm -rf src; fi
	mkdir -p src/github.com/whosonfirst/go-whosonfirst-markdown
	cp -r flags src/github.com/whosonfirst/go-whosonfirst-markdown/
	cp -r jekyll src/github.com/whosonfirst/go-whosonfirst-markdown/
	cp -r parser src/github.com/whosonfirst/go-whosonfirst-markdown/
	cp -r render src/github.com/whosonfirst/go-whosonfirst-markdown/
	cp -r search src/github.com/whosonfirst/go-whosonfirst-markdown/
	cp -r utils src/github.com/whosonfirst/go-whosonfirst-markdown/
	cp -r writer src/github.com/whosonfirst/go-whosonfirst-markdown/
	cp *.go src/github.com/whosonfirst/go-whosonfirst-markdown/
	cp -r vendor/* src/

rmdeps:
	if test -d src; then rm -rf src; fi 

build:	fmt bin

deps:   
	@GOPATH=$(GOPATH) go get -u "github.com/microcosm-cc/bluemonday"
	@GOPATH=$(GOPATH) go get -u "gopkg.in/russross/blackfriday.v2"
	@GOPATH=$(GOPATH) go get -u "github.com/djherbis/times"
	@GOPATH=$(GOPATH) go get -u "github.com/facebookgo/atomicfile"
	@GOPATH=$(GOPATH) go get -u "github.com/whosonfirst/go-whosonfirst-crawl"

vendor-deps: rmdeps deps
	if test ! -d vendor; then mkdir vendor; fi
	if test -d vendor; then rm -rf vendor; fi
	cp -r src vendor
	find vendor -name '.git' -print -type d -exec rm -rf {} +
	rm -rf src

fmt:
	go fmt cmd/*.go
	go fmt flags/*.go
	go fmt jekyll/*.go
	go fmt parser/*.go
	go fmt render/*.go
	go fmt search/*.go
	go fmt utils/*.go
	go fmt writer/*.go
	go fmt *.go

bin: 	rmdeps self
	rm -rf bin/*
	@GOPATH=$(shell pwd) go build -o bin/wof-markdown-parse cmd/wof-markdown-parse.go
	@GOPATH=$(shell pwd) go build -o bin/wof-md2feed cmd/wof-md2feed.go	
	@GOPATH=$(shell pwd) go build -o bin/wof-md2html cmd/wof-md2html.go
	@GOPATH=$(shell pwd) go build -o bin/wof-md2idx cmd/wof-md2idx.go
