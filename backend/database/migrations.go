package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"fmt"
	"slices"
)

func RunMigrations(db *sql.DB, populateDb bool) error {
	cwd, _ := os.Getwd()
	fmt.Println(cwd)
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return err
	}
	log.Println("FOUND MIGRATIONS:", files)

	slices.Sort(files)

	if populateDb {
	log.Println("NO FILL DB")
	files = files[:5]	
	}

	for _, file := range files {
		log.Println("running migrations", file)
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
