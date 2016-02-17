package storage

import (
	"database/sql"
)

const entryTABLE = `CREATE TABLE "log" (
		"id" int(11) NOT NULL,
		"host" text NOT NULL,
		"port" int(2) NOT NULL,
		"insertTime" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		"activeTime" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`

func openDb(dbName string) (*sql.DB, error) {
	db, err := sql.Open(dbName, "sqlite3")
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	tx.Exec(entryTABLE)

	tx.Commit()

	return db, err
}
