package tables

/*

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

*/
