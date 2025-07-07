package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
)

func Migrate() {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		filename VARCHAR(255) UNIQUE
	)`)
	if err != nil {
		panic(err)
	}

	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		panic(err)
	}

	sort.Strings(files)

	for _, file := range files {
		filename := filepath.Base(file)

		var exists int
		err := DB.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = ?", filename).Scan(&exists)
		if err != nil {
			panic(err)
		}

		if exists == 0 {
			content, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}

			_, err = DB.Exec(string(content))
			if err != nil {
				panic(err)
			}

			_, err = DB.Exec("INSERT INTO migrations (filename) VALUES (?)", filename)
			if err != nil {
				panic(err)
			}

			fmt.Println("✅ Migrated:", filename)
		} else {
			fmt.Println("⏩ Skipped:", filename)
		}
	}
}
