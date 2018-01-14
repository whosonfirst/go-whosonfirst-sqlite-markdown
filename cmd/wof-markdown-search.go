package main

import (
	"flag"
	"github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-markdown/parser"
	"github.com/whosonfirst/go-whosonfirst-markdown-search"
	"log"
)

func main() {

	var index = flag.String("index", "bleve", "")
	var dsn = flag.String("dsn", "index.db", "")
	var q = flag.String("query", "", "")

	flag.Parse()

	var idx search.Indexer

	switch *index {

	case "bleve":

		i, err := search.NewBleveIndexer(*dsn)

		if err != nil {
			log.Fatal(err)
		}

		idx = i

	case "sqlite":

		i, err := search.NewSQLiteIndexer(*dsn)

		if err != nil {
			log.Fatal(err)
		}

		idx = i
	default:
		log.Fatal("Invalid indexer")
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
