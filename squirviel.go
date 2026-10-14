package main

import (
	"database/sql"
	"fmt"

	sqlite "gosqlite.org"
)

var db = setupDb()

func setupDb() *sqlite.DB {
	db, err := sqlite.Open(sqlite.Config{
		Path:         "hose_paths.db",
		Pragmas:      sqlite.RecommendedPragmas(), // WAL + busy_timeout=5s + foreign_keys
		MaxOpenConns: 8,
	})

	if err != nil {
		fmt.Print(fmt.Errorf("open database: %w", err))
	}

	createQuery := `CREATE TABLE IF NOT EXISTS projects (
		name STRING NOT NULL PRIMARY KEY,
		rootPath STRING NOT NULL,
		keyStorePath STRING NOT NULL,
		alias STRING NOT NULL
	)`
	_, errr := db.Exec(createQuery)

	if err != nil || errr != nil {
		fmt.Print(fmt.Errorf("open database: %w", err))
	}
	return db
}

func getAllEntries() *sql.Rows {
	rows, _ := db.Query("SELECT * FROM projects")
	return rows
}

func insertProject(p project) sql.Result {
	query := "INSERT INTO projects(name, rootPath, keyStorePath, alias) VALUES(?, ?, ?, ?)"
	res, err := db.Exec(query, p.name, p.rootPath, p.keyStorePath, p.alias)
	if err != nil {

		fmt.Print(fmt.Errorf("insert: %w", err))
	}
	return res
}
