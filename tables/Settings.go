package tables

import (
	"sync"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/mohammedaouamri5/JDM-back/db"
	"github.com/sirupsen/logrus"
)

type _Settings struct {
	Working_dir string
	Output_dir  string
}

var (
	instance *_Settings
	once     sync.Once
)

func select_Settings() {
	instance = &_Settings{}

	sql, _, err := sq.Select("*").From("settings").ToSql()
	if err != nil {
		logrus.Fatal(err)
	}

	rows, err := DB().Query(sql)
	if err != nil {
		logrus.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&instance.Working_dir,
			&instance.Output_dir,
		)
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func Settings() *_Settings {
	once.Do(select_Settings)
	return instance
}
