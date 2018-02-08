# go-whosonfirst-markdown

There are many Markdown tools. This one is ours.

## Install

You will need to have both `Go` (specifically a version of Go more recent than 1.7 so let's just assume you need [Go 1.9](https://golang.org/dl/) or higher) and the `make` programs installed on your computer. Assuming you do just type:

```
make bin
```

All of this package's dependencies are bundled with the code in the `vendor` directory.

## Important

Everything is in flux, right now. Lots of things will change.

## Tools

### wof-md2html

```
./bin/wof-md2html -h
Usage of ./bin/wof-md2html:
  -footer string
    	The name of the (Go) template to use as a custom footer
  -header string
    	The name of the (Go) template to use as a custom header
  -input string
    	What you expect the input Markdown file to be called (default "index.md")
  -mode string
    	Valid modes are: files, directory (default "files")
  -output string
    	What you expect the output HTML file to be called (default "index.html")
  -templates value
    	One or more templates to parse in addition to -header and -footer
  -writer value
    	One or more writer to output rendered Markdown to. Valid writers are: fs=PATH; null; stdout
```

### wof-md2idx

```
./bin/wof-md2idx -h
Usage of ./bin/wof-md2idx:
  -footer string
    	The name of the (Go) template to use as a custom footer
  -header string
    	The name of the (Go) template to use as a custom header
  -input string
    	What you expect the input Markdown file to be called (default "index.md")
  -output string
    	What you expect the output HTML file to be called (default "index.html")
  -templates value
    	One or more templates to parse in addition to -header and -footer
  -writer value
    	One or more writer to output rendered Markdown to. Valid writers are: fs=PATH; null; stdout
```

### wof-md2feed

```
./bin/wof-md2feed -h
Usage of ./bin/wof-md2feed:
  -format string
    	Valid options are: atom_10, rss_20 (default "rss_20")
  -input string
    	What you expect the input Markdown file to be called (default "index.md")
  -items int
    	The number of items to include in your feed (default 10)
  -output string
    	The filename of your feed. If empty default to the value of -format + ".xml"
  -templates value
    	One or more directories containing (Go) templates to parse
  -writer value
    	One or more writer to output rendered Markdown to. Valid writers are: fs=PATH; null; stdout
```

### wof-mdparse

```
./bin/wof-mdparse -h
Usage of ./bin/wof-mdparse:
  -all
    	Dump both frontmatter and body
  -body
    	Dump (Markdown) body
  -frontmatter
    	Dump (Jekyll) frontmatter
```

## See also

* github.com/microcosm-cc/bluemonday
* gopkg.in/russross/blackfriday.v2
* https://github.com/whosonfirst/go-whosonfirst-markdown-sqlite
* https://github.com/whosonfirst/go-whosonfirst-markdown-bleve
