package main

import (
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-sqlite/database"
	"github.com/whosonfirst/go-whosonfirst-sqlite/utils"
	"io"
	"os"
	"strings"
)

func main() {

	driver := flag.String("driver", "sqlite3", "")
	var dsn = flag.String("dsn", "index.db", "")

	var table = flag.String("table", "documents_search", "")

	flag.Parse()

	logger := log.SimpleWOFLogger()

	stdout := io.Writer(os.Stdout)
	logger.AddLogger(stdout, "status")

	db, err := database.NewDBWithDriver(*driver, *dsn)

	if err != nil {
		logger.Fatal("unable to create database (%s) because %s", *dsn, err)
	}

	defer db.Close()

	ok, err := utils.HasTable(db, *table)

	if err != nil {
		logger.Fatal("unable to test for table (%s) because %s", *table, err)
	}

	if !ok {
		logger.Fatal("database is missing table (%s)", *table)
	}

	conn, err := db.Conn()

	if err != nil {
		logger.Fatal("CONN", err)
	}

	q := strings.Join(flag.Args(), " ")

	sql := fmt.Sprintf("SELECT id FROM %s WHERE body MATCH ?", *table)

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
