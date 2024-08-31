package main

import (
	"database/sql"
	"log"
	"os"
)


func InitTables(db *sql.DB) error {

	files := []string{
		"./tables/tables.sql",
		"./tables/fill.sql", 
	}

	for _, file := range files {
		sqlFile, err := os.ReadFile(file)
		
		if err != nil {
			log.Fatal(err)
		} else if _, err := db.Exec(string(sqlFile)); err != nil {
			log.Fatal(err)
			return err
		} else {
			log.Println("init", file)
		}
	}

	return nil
}
