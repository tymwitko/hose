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

func getAllEntries() ([]project, error) {
	rows, err := db.Query(`SELECT * FROM projects`)
	if err != nil {
		return nil, fmt.Errorf("query projects: %w", err)
	}
	defer rows.Close()

	var projects []project
	for rows.Next() {
		var p project
		if err := rows.Scan(&p.name, &p.rootPath, &p.keyStorePath, &p.alias); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}
	return projects, nil
}

func insertProject(p project) sql.Result {
	query := "INSERT INTO projects(name, rootPath, keyStorePath, alias) VALUES(?, ?, ?, ?)"
	res, err := db.Exec(query, p.name, p.rootPath, p.keyStorePath, p.alias)
	if err != nil {

		fmt.Print(fmt.Errorf("insert: %w", err))
	}
	return res
}
