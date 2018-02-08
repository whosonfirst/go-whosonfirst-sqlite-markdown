package tables

import (
	"fmt"
	md "github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-markdown/search"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite-markdown"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
)

type DocumentsLinksTable struct {
	markdown.MarkdownTable
	name string
}

func NewDocumentsLinksTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewDocumentsLinksTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewDocumentsLinksTable() (sqlite.Table, error) {

	t := DocumentsLinksTable{
		name: "documents_links",
	}

	return &t, nil
}

func (t *DocumentsLinksTable) Name() string {
	return t.name
}

func (t *DocumentsLinksTable) Schema() string {

	schema := `CREATE TABLE %s (
	       document_id TEXT,
	       host TEXT,
	       link TEXT
	);

	CREATE UNIQUE INDEX %s_by_link ON %s (document_id, host, link);`

	// this is a bit stupid really... (20170901/thisisaaronland)
	return fmt.Sprintf(schema, t.Name(), t.Name(), t.Name())
}

func (t *DocumentsLinksTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *DocumentsLinksTable) IndexRecord(db sqlite.Database, i interface{}) error {
	return t.IndexDocument(db, i.(*md.Document))
}

func (t *DocumentsLinksTable) IndexDocument(db sqlite.Database, doc *md.Document) error {

	search_doc, err := search.NewSearchDocument(doc)

	if err != nil {
		return err
	}

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	// DELETE FROM documents_links WHERE document_id = ?

	document_id := search_doc.Id

	// maybe dissolve links into their own table not associated with documents
	// and simply store post_id, link_id... in a `documents_links` table - maybe
	// but not today (2018011/thisisaaronland)

	for _, url := range search_doc.Links {

		sql := fmt.Sprintf(`INSERT OR REPLACE INTO documents_links (document_id, host, link) VALUES (?, ?, ?)`)
		stmt, err := tx.Prepare(sql)

		if err != nil {
			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(document_id, url.Host, url.Path)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
