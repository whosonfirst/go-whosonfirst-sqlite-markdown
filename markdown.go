package markdown

import (
	md "github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
)

type DocumentsTable interface {
	sqlite.Table
	IndexDocument(sqlite.Database, md.Document) error
}
