package tables

import (
	"fmt"
	md "github.com/whosonfirst/go-whosonfirst-markdown"
	"github.com/whosonfirst/go-whosonfirst-markdown/search"
	"github.com/whosonfirst/go-whosonfirst-sqlite"
	"github.com/whosonfirst/go-whosonfirst-sqlite-markdown"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
)

type DocumentsTables struct {
	markdown.DocumentsTable
	name string
}

func NewDocumentsTableWithDatabase(db sqlite.Database) (sqlite.Table, error) {

	t, err := NewDocumentsTable()

	if err != nil {
		return nil, err
	}

	err = t.InitializeTable(db)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func NewDocumentsTable() (sqlite.Table, error) {

	t := DocumentsTable{
		name: "documents",
	}

	return &t, nil
}

func (t *DocumentsTable) Name() string {
	return t.name
}

func (t *DocumentsTable) Schema() string {

	schema := `CREATE TABLE %s (
	       id TEXT PRIMARY KEY,
	       title TEXT,
	       category TEXT,
	       date TEXT,
	       body TEXT,
	       code TEXT
	);

	CREATE INDEX documents_by_date ON %s (date);
	CREATE INDEX documents_by_body ON %s (body);
	`

	// this is a bit stupid really... (20170901/thisisaaronland)
	return fmt.Sprintf(sql, t.Name(), t.Name(), t.Name(), t.Name())
}

func (t *DocumentsTable) InitializeTable(db sqlite.Database) error {

	return utils.CreateTableIfNecessary(db, t)
}

func (t *DocumentsTable) IndexRecord(db sqlite.Database, i interface{}) error {
	return t.IndexDocument(db, i.(md.Document))
}

func (t *DocumentsTable) IndexDocument(db sqlite.Database, doc md.Document) error {

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

	str_body := strings.Join(search_doc.Body, " ")
	str_code := strings.Join(search_doc.Code, " ")

	sql := fmt.Sprintf(`INSERT OR REPLACE INTO documents (
			id, title, category, date, body, code
		) VALUES (
			?, ?, ?, ?, ?, ?
		)`)

	stmt, err := tx.Prepare(sql)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(search_doc.Id, search_doc.Title, search_doc.Category, search_doc.Date, str_body, str_code)

	if err != nil {
		return err
	}

	return tx.Commit()
}
