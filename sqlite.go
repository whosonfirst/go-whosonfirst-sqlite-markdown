package search

// https://sqlite.org/fts3.html
// https://www.sqlite.org/fts5.html

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/whosonfirst/go-whosonfirst-markdown"
	_ "log"
	"strings"
)

type SQLiteIndexer struct {
	Indexer
	conn *sql.DB
	dsn  string
}

func NewSQLiteIndexer(dsn string) (Indexer, error) {

	// It seems likely to me that once we understand how this works
	// we will replace most of the code below with go-whosonfirst-sqlite
	// and a series of markdown specific tables but not today...
	// (20180110/thisisaaronland)

	// or maybe not until the interface for go-whosonfirst-sqlite is
	// updated to index interface{} rather than geojson.Feature - really
	// I don't know yet... (20180110/thisisaaronland)

	conn, err := sql.Open("sqlite3", dsn)

	if err != nil {
		return nil, err
	}

	sql := "SELECT name FROM sqlite_master WHERE type='table'"

	rows, err := conn.Query(sql)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	has_table := false

	for rows.Next() {

		var name string
		err := rows.Scan(&name)

		if err != nil {
			return nil, err
		}

		if name == "documents" {
			has_table = true
			break
		}
	}

	if !has_table {

		// this needs a "tags" table but don't bother adding that until
		// we figure out what do about using or not using go-whosonfirst-sqlite
		// above (20180110/thisisaaronland)

		schema := `CREATE TABLE documents (
		       id TEXT PRIMARY KEY,
		       title TEXT,
		       category TEXT,
		       date TEXT,
		       body TEXT,
		       code TEXT
		);

		CREATE INDEX documents_by_date ON documents (date);
		CREATE INDEX documents_by_body ON documents (body);

		CREATE VIRTUAL TABLE documents_search USING fts4(id, title, category, body, code);

		CREATE TABLE documents_authors (
		       document_id TEXT,
		       author TEXT,
		       date TEXT
		);

		CREATE UNIQUE INDEX documents_authors_by_author ON documents_authors (document_id, author);
		CREATE INDEX documents_authors_by_date ON documents_authors (author, date);

		CREATE TABLE documents_links (
		       document_id TEXT,
		       host TEXT,
		       link TEXT
		);

		CREATE UNIQUE INDEX documents_links_by_link ON documents_links (document_id, host, link);
		`
		_, err = conn.Exec(schema)

		if err != nil {
			return nil, err
		}
	}

	i := SQLiteIndexer{
		conn: conn,
		dsn:  dsn,
	}

	return &i, nil
}

func (i *SQLiteIndexer) Conn() (*sql.DB, error) {
	return i.conn, nil
}

func (i *SQLiteIndexer) DSN() string {
	return i.dsn
}

func (i *SQLiteIndexer) Close() error {
	return i.conn.Close()
}

func (i *SQLiteIndexer) Query(q *SearchQuery) (interface{}, error) {

	conn, err := i.Conn()

	if err != nil {
		return nil, err
	}

	sql := fmt.Sprintf("SELECT * FROM documents_search(?)")

	rows, err := conn.Query(sql, q.QueryString)

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (i *SQLiteIndexer) IndexDocument(doc *markdown.Document) (*SearchDocument, error) {

	search_doc, err := NewSearchDocument(doc)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = i.IndexDocumentsTable(ctx, search_doc)

	if err != nil {
		return nil, err
	}

	err = i.IndexDocumentsAuthorsTable(ctx, search_doc)

	if err != nil {
		return nil, err
	}

	err = i.IndexDocumentsLinksTable(ctx, search_doc)

	if err != nil {
		return nil, err
	}

	err = i.IndexDocumentsSearchTable(ctx, search_doc)

	if err != nil {
		return nil, err
	}

	return search_doc, nil
}

func (i *SQLiteIndexer) IndexDocumentsTable(ctx context.Context, search_doc *SearchDocument) error {

	select {

	case <-ctx.Done():
		return nil
	default:

		conn, err := i.Conn()

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

		err = tx.Commit()

		if err != nil {
			return err
		}

		return nil
	}
}

func (i *SQLiteIndexer) IndexDocumentsAuthorsTable(ctx context.Context, search_doc *SearchDocument) error {

	select {

	case <-ctx.Done():
		return nil
	default:

		conn, err := i.Conn()

		if err != nil {
			return err
		}

		tx, err := conn.Begin()

		if err != nil {
			return err
		}

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

		err = tx.Commit()

		if err != nil {
			return err
		}

		return nil
	}
}

func (i *SQLiteIndexer) IndexDocumentsLinksTable(ctx context.Context, search_doc *SearchDocument) error {

	select {

	case <-ctx.Done():
		return nil
	default:

		conn, err := i.Conn()

		if err != nil {
			return err
		}

		tx, err := conn.Begin()

		if err != nil {
			return err
		}

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

		err = tx.Commit()

		if err != nil {
			return err
		}

		return nil
	}
}

func (i *SQLiteIndexer) IndexDocumentsSearchTable(ctx context.Context, search_doc *SearchDocument) error {

	select {

	case <-ctx.Done():
		return nil
	default:

		conn, err := i.Conn()

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

		err = tx.Commit()

		if err != nil {
			return err
		}

		return nil
	}
}
