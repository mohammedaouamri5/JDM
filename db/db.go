package db

import (
	"database/sql"
	"errors" 
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
	"github.com/sirupsen/logrus"
)

type _DB struct {
	db *sql.DB
}

var (
	instance *_DB
	once     sync.Once
)

func initTables() error {
	if instance == nil {
		return errors.New("db not initialized")
	}

	files := []string{
		"./tables/tables.sql",
		"./tables/fill.sql",
	}

	for _, file := range files {
		sqlFile, err := os.ReadFile(file)
		if err != nil {
			logrus.Error(err)
			return err
		}
		if _, err := instance.db.Exec(string(sqlFile)); err != nil {
			logrus.Error(err)
			return err
		} else {
			logrus.Info("Initialized: ", file)
		}
	}

	return nil
}

func DB() *sql.DB {
	once.Do(func() {
		instance = &_DB{}
		var err error
		instance.db, err = sql.Open("sqlite3", "./tmp/jdm.db")
		if err != nil {
			logrus.Fatal(err)
		}

		if err := initTables(); err != nil {
			logrus.Fatal(err)
		}
	})

	return instance.db
}
