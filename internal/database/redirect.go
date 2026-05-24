package database

import (
	"database/sql"

	"github.com/beecorrea/shortlinks/internal/structs"
	_ "modernc.org/sqlite"
)

type RedirectDB struct {
	db *sql.DB
}

func NewRedirectDB(dbPath string) (*RedirectDB, error) {
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
	if err != nil {
		return err
	}
	_, _ = db.Exec("ALTER TABLE links ADD COLUMN summary TEXT")
	return nil
}

func (r *RedirectDB) CreateShortlink(key, url, domain, summary string) error {
	_, err := r.db.Exec(queryInsertShortlink, key, url, domain, summary)
	return err
}

func (r *RedirectDB) GetShortlinkByKey(key string) (*structs.Shortlink, error) {
	row := r.db.QueryRow(queryGetShortlinkByKey, key)

	var s structs.Shortlink
	err := row.Scan(&s.Key, &s.URL, &s.Domain, &s.Summary, &s.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *RedirectDB) ListShortlinks() ([]structs.Shortlink, error) {
	rows, err := r.db.Query(queryListShortlinks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shortlinks []structs.Shortlink
	for rows.Next() {
		var s structs.Shortlink
		if err := rows.Scan(&s.Key, &s.URL, &s.Domain, &s.Summary, &s.CreatedAt); err != nil {
			return nil, err
		}
		shortlinks = append(shortlinks, s)
	}
	return shortlinks, nil
}

func (r *RedirectDB) DeleteShortlink(key string) error {
	_, err := r.db.Exec(queryDeleteShortlink, key)
	return err
}
