package search

import (
	"github.com/blevesearch/bleve"
	"github.com/whosonfirst/go-whosonfirst-markdown"
	_ "log"
	"os"
)

type BleveIndexer struct {
	Indexer
	index bleve.Index
}

func NewBleveIndexer(path string) (Indexer, error) {

	var idx bleve.Index

	_, err := os.Stat(path)

	if os.IsNotExist(err) {

		m := bleve.NewIndexMapping()
		b, err := bleve.New(path, m)

		if err != nil {
			return nil, err
		}

		idx = b
	} else {

		b, err := bleve.Open(path)

		if err != nil {
			return nil, err
		}

		idx = b
	}

	i := BleveIndexer{
		index: idx,
	}

	return &i, err
}

func (i *BleveIndexer) Close() error {
	return nil
}

func (i *BleveIndexer) Query(q *SearchQuery) (interface{}, error) {
	query := bleve.NewQueryStringQuery(q.QueryString)
	req := bleve.NewSearchRequest(query)
	return i.index.Search(req)
}

func (i *BleveIndexer) IndexDocument(doc *markdown.Document) (*SearchDocument, error) {

	search_doc, err := NewSearchDocument(doc)

	if err != nil {
		return nil, err
	}

	err = i.index.Index(search_doc.Id, search_doc)

	if err != nil {
		return nil, err
	}

	return search_doc, nil
}
