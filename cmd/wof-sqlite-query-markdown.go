package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"io"
	"os"
	"strings"
)

func main() {

	driver := flag.String("driver", "sqlite3", "")
	var dsn = flag.String("dsn", "index.db", "")

	var table = flag.String("table", "documents_search", "")
	var col = flag.String("column", "body", "")

	flag.Parse()

	logger := log.SimpleWOFLogger()

	stdout := io.Writer(os.Stdout)
	logger.AddLogger(stdout, "status")

	db, err := database.NewDBWithDriver(*driver, *dsn)

	if err != nil {
		logger.Fatal("unable to create database (%s) because %s", *dsn, err)
	}

	defer db.Close()

	conn, err := db.Conn()

	if err != nil {
		logger.Fatal("CONN", err)
	}

	q := strings.Join(flag.Args(), " ")

	sql := fmt.Sprintf("SELECT id FROM %s WHERE %s MATCH ?", *table, *col)
	rows, err := conn.Query(sql, q)

	if err != nil {
		logger.Fatal("QUERY", err)
	}

	defer rows.Close()

	for rows.Next() {

		var id string

		err = rows.Scan(&id)

		if err != nil {
			logger.Fatal("ID", err)
		}

		logger.Status("%s - %s", q, id)
	}

	err = rows.Err()

	if err != nil {
		logger.Fatal("ROWS", err)
	}

}
