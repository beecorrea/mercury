package database

import (
	"database/sql"

	"github.com/beecorrea/shortlinks/internal/structs"
	_ "modernc.org/sqlite"
)

type RedirectDB struct {
	db *sql.DB
}

func New(dbPath string) (*RedirectDB, error) {
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

	return &RedirectDB{db: db}, nil
}

func (r *RedirectDB) Close() error {
	return r.db.Close()
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(queryMigrate)
	return err
}

func (r *RedirectDB) CreateLink(key, url, domain string) error {
	_, err := r.db.Exec(queryInsertLink, key, url, domain)
	return err
}

func (r *RedirectDB) GetLinkByKey(key string) (*structs.Link, error) {
	row := r.db.QueryRow(queryGetLinkByKey, key)

	var link structs.Link
	err := row.Scan(&link.Key, &link.URL, &link.Domain, &link.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *RedirectDB) ListLinks() ([]structs.Link, error) {
	rows, err := r.db.Query(queryListLinks)
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

func (r *RedirectDB) DeleteLink(key string) error {
	_, err := r.db.Exec(queryDeleteLink, key)
	return err
}
