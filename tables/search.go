package tables

import (
	"fmt"
	md "github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-markdown/search"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite-markdown"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
	"strings"
)

type DocumentsSearchTable struct {
	markdown.MarkdownTable
	name string
}

func NewDocumentsSearchTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewDocumentsSearchTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewDocumentsSearchTable() (sqlite.Table, error) {

	t := DocumentsSearchTable{
		name: "documents_search",
	}

	return &t, nil
}

func (t *DocumentsSearchTable) Name() string {
	return t.name
}

func (t *DocumentsSearchTable) Schema() string {

	schema := `CREATE VIRTUAL TABLE %s USING fts4(id, title, category, body, code);`

	// this is a bit stupid really... (20170901/thisisaaronland)
	return fmt.Sprintf(schema, t.Name())
}

func (t *DocumentsSearchTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *DocumentsSearchTable) IndexRecord(db sqlite.Database, i interface{}) error {
	return t.IndexDocument(db, i.(*md.Document))
}

func (t *DocumentsSearchTable) IndexDocument(db sqlite.Database, doc *md.Document) error {

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

	str_body := strings.Join(search_doc.Body, " ")
	str_code := strings.Join(search_doc.Code, " ")

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO documents_search (
		id, title, category, body, code
			) VALUES (
		?, ?, ?, ?, ?
		)`)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(search_doc.Id, search_doc.Title, search_doc.Category, str_body, str_code)

	if err != nil {
		return err
	}

	return tx.Commit()
}
