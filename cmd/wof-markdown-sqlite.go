package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-markdown/parser"
	"github.com/whosonfirst/go-whosonfirst-markdown/search"
	"github.com/whosonfirst/go-whosonfirst-markdown-sqlite"
	"log"
)

func main() {

	var dsn = flag.String("dsn", "index.db", "")
	var q = flag.String("query", "", "")

	flag.Parse()

		idx, err := sqlite.NewSQLiteIndexer(*dsn)

		if err != nil {
			log.Fatal(err)
		}

	opts := parser.DefaultParseOptions()

	for _, path := range flag.Args() {

		fm, b, err := parser.ParseFile(path, opts)

		if err != nil {
			log.Fatal(err)
		}

		doc, err := markdown.NewDocument(fm, b)

		if err != nil {
			log.Fatal(err)
		}

		search_doc, err := idx.IndexDocument(doc)

		if err != nil {
			log.Fatal(err)
		}

		log.Println("INDEX", path, search_doc.Title)
	}

	if *q != "" {

		query, err := search.NewDefaultSearchQuery(*q)

		if err != nil {
			log.Fatal(err)
		}

		r, err := idx.Query(query)

		if err != nil {
			log.Fatal(err)
		}

		log.Println(r)
	}
}
