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
	queryMigrate      = readQuery("migrate")
	queryInsertLink   = readQuery("insert_link")
	queryGetLinkByKey = readQuery("get_link")
	queryListLinks    = readQuery("list_links")
	queryDeleteLink   = readQuery("delete_link")
)
