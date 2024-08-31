package main

import (
	"database/sql"
	"log" 

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./tmp/jdm.db")
	if err != nil {
		log.Fatal(err)
	}

	if err := InitTables(db) ; err != nil {
		log.Fatal(err)
	}

	


}