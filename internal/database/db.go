package database

import (
	"database/sql"

	"github.com/beecorrea/shortlinks/internal/structs"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(queryMigrate)
	return err
}

func (db *DB) CreateLink(key, url, domain string) error {
	_, err := db.Exec(queryInsertLink, key, url, domain)
	return err
}

func (db *DB) GetLinkByKey(key string) (*structs.Link, error) {
	row := db.QueryRow(queryGetLinkByKey, key)

	var link structs.Link
	err := row.Scan(&link.Key, &link.URL, &link.Domain, &link.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &link, nil
}

func (db *DB) ListLinks() ([]structs.Link, error) {
	rows, err := db.Query(queryListLinks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []structs.Link
	for rows.Next() {
		var link structs.Link
		if err := rows.Scan(&link.Key, &link.URL, &link.Domain, &link.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (db *DB) DeleteLink(key string) error {
	_, err := db.Exec(queryDeleteLink, key)
	return err
}
