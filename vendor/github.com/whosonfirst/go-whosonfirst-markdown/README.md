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

## Example

```
./bin/wof-md2html -writer fs=../whosonfirst-www/www/blog -header ../whosonfirst-www/templates/blog/header.html -footer ../whosonfirst-www/templates/blog/footer.html -mode directory ../whosonfirst-www/blog/
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/06/24/sf-neighbourhood-updates/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/02/19/iamhere/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/04/14/missing-the-point/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2015/09/28/spelunker-jumping-into-who-s-on-first/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/07/13/wikipedia-data/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2015/08/18/who-s-on-first/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/10/06/wof-lifecycle-documentation/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/10/05/boundary-issues-properties/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/12/08/mesoshapes/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/10/07/whosonfirst-venues/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/04/08/yesnofix/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/08/24/all-of-the-places/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2016/08/15/mapping-with-bias/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/04/20/neighbourhood-updates-two/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/04/04/whosonfirst-api/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/05/10/simple-is-hard/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/04/17/weird-and-wonderful/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/02/10/bundler/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/08/01/geotagging-wof-venues/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/09/19/introducting-statoids/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/06/29/tackling-space-and-time-in-whosonfirst/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/08/22/summer-2017-wof/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/10/17/whosonfirst-nacis-2017/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/12/14/updating-whosonfirst/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/04/28/whosonfirst-api-explorer/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/10/24/whosonfirst-sotmus-2017/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/12/22/neighbourhood-updates-three/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/10/05/mapzen-places/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/07/28/wof-website-redesign/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/12/21/wof-in-a-box/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2018/01/02/chapter-two/index.html
2017/12/30 17:07:30 wrote /usr/local/mapzen/whosonfirst-www/content/blog/2017/12/29/wof-in-a-box-part2/index.html
```

## Assumptions

1. That your Markdown files have Jekyll-style "front matter"
2. That you have one Markdown file per directory

## What is "front matter"

"Front matter" is metadata included between `---` brackets at the top of a Markdown file. For example:

```
---
layout: page
category: blog
title: Bundling up descendants into GeoJSON
excerpt: We made a handy tool that lets you download the descendants of a place as GeoJSON.
image: images/wof_canada-descender.gif
authors: [burritojustice, stepps00, dphiffer]
tags: [whosonfirst, data]
---

Every record in [Who’s On First](https://whosonfirst.mapzen.com/spelunker/), our gazetteer of places around the world, has two powerful properties -- geometry and hierarchy. Since the raw data of any WOF record is a blob of GeoJSON, you can take the URL and draw it in any modern mapping application with minimal effort. We thought it would be nice to let you grab GeoJSON for a bunch of records at once, so we added a tool called the Bundler for you using the soon-to-be-released Who’s On First API.
```

## Tools

### wof-md2html

Render one or more Markdown files as HTML.

```
./bin/wof-md2html -h
Usage of ./bin/wof-md2html:
  -footer string
    	The path to a custom (Go) template to use as a footer for your HTML output
  -header string
    	The path to a custom (Go) template to use as header for your HTML output
  -input string
    	What you expect the input Markdown file to be called (default "index.md")
  -mode string
    	Valid modes are: files, directory (default "files")
  -output string
    	What you expect the output HTML file to be called (default "index.html")
  -writer value
    	One or more writer to output rendered Markdown to. Valid writers are: fs=PATH; null; stdout
```

### wof-markdown-to-search

_This tool will be renamed. Also it is incomplete. Hence the weird results below, but you get the idea..._

```
./bin/wof-markdown-search -query bummer examples/index.md 
2018/01/07 18:37:25 INDEX examples/index.md Who's On First, Chapter Two
2018/01/07 18:37:25 2 matches, showing 1 through 2, took 251.079µs
    1. debug (0.008869)
    2. Who's On First, Chapter Two (0.008866)
```

Currently this uses [Bleve](https://github.com/blevesearch/bleve) under the hood.

## See also

* github.com/microcosm-cc/bluemonday
* gopkg.in/russross/blackfriday.v2
