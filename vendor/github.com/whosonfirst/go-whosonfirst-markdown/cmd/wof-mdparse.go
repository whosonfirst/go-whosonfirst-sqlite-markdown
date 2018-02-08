package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-markdown/parser"
	"log"
)

func main() {

	var frontmatter = flag.Bool("frontmatter", false, "Dump (Jekyll) frontmatter")
	var body = flag.Bool("body", false, "Dump (Markdown) body")
	var all = flag.Bool("all", false, "Dump both frontmatter and body")

	flag.Parse()

	if *all {
		*frontmatter = true
		*body = true
	}

	opts := parser.DefaultParseOptions()
	opts.FrontMatter = *frontmatter
	opts.Body = *body

	for _, path := range flag.Args() {

		fm, b, err := parser.ParseFile(path, opts)

		if err != nil {
			log.Fatal(err)
		}

		if *frontmatter {
			fmt.Println(fm.String())
		}

		if *body {
			fmt.Println(b.String())
		}

	}
}
