package tables

import (
	"fmt"
	md "github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-markdown/search"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite-markdown"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
)

type DocumentsAuthorsTables struct {
	markdown.DocumentsAuthorsTable
	name string
}

func NewDocumentsAuthorsTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewDocumentsAuthorsTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewDocumentsAuthorsTable() (sqlite.Table, error) {

	t := DocumentsAuthorsTable{
		name: "documents_authors",
	}

	return &t, nil
}

func (t *DocumentsAuthorsTable) Name() string {
	return t.name
}

func (t *DocumentsAuthorsTable) Schema() string {

	schema := `CREATE TABLE %s (
 	       document_id TEXT,
	       author TEXT,
	       date TEXT
	);

	CREATE UNIQUE INDEX %s_by_author ON %s (document_id, author);
	CREATE INDEX %s_by_date ON %s (author, date);	`

	// this is a bit stupid really... (20170901/thisisaaronland)
	return fmt.Sprintf(sql, t.Name(), t.Name(), t.Name(), t.Name(), t.Name())
}

func (t *DocumentsAuthorsTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *DocumentsAuthorsTable) IndexRecord(db sqlite.Database, i interface{}) error {
	return t.IndexDocument(db, i.(md.Document))
}

func (t *DocumentsAuthorsTable) IndexDocument(db sqlite.Database, doc md.Document) error {

	search_doc, err := search.NewSearchDocument(doc)

	if err != nil {
		return nil, err
	}

	conn, err := db.Conn()

	if err != nil {
		return err
	}

	tx, err := conn.Begin()

	if err != nil {
		return err
	}

	// DELETE FROM documents_authors WHERE document_id = ?

	document_id := search_doc.Id
	date := search_doc.Date

	for _, author := range search_doc.Authors {

		sql := fmt.Sprintf(`INSERT OR REPLACE INTO documents_authors (document_id, author, date) VALUES (?, ?, ?)`)
		stmt, err := tx.Prepare(sql)

		if err != nil {
			return err
		}

		defer stmt.Close()

		_, err = stmt.Exec(document_id, author, date)

		if err != nil {
			return err
		}

	}

	return tx.Commit()
}
