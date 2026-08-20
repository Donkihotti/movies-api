package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"slices"
)

func RunMigrations(db *sql.DB) error {

	files, err := filepath.Glob("database/migrations/*.sql")
	if err != nil {
		return err
	}

	slices.Sort(files)
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(content))
		if err != nil {
			return err
		}
	}
	return nil
}
