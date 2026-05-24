package database

import (
	"embed"
	"strings"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

// readQuery reads the query file from the embedded filesystem.
func readQuery(name string) string {
	content, err := sqlFiles.ReadFile("sql/" + name + ".sql")
	if err != nil {
		panic("failed to read embedded SQL query " + name + ": " + err.Error())
	}
	return strings.TrimSpace(string(content))
}

var (
	queryMigrate                 = readQuery("migrate")
	queryInsertShortlink         = readQuery("insert_link")
	queryGetShortlinkByKey       = readQuery("get_link")
	queryListShortlinks          = readQuery("list_links")
	queryDeleteShortlink         = readQuery("delete_link")
	queryListFailedScrapes       = readQuery("list_failed_scrapes")
	queryUpdateSummary           = readQuery("update_summary")
	queryIncrementScrapeAttempts = readQuery("increment_scrape_attempts")
)

