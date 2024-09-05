package main

import (
	_ "github.com/mattn/go-sqlite3"
	. "github.com/mohammedaouamri5/JDM-back/db"
	"github.com/mohammedaouamri5/JDM-back/tables"
	"github.com/sirupsen/logrus"
)

func main() {

	InitLog() // No Error handling yet
	DB()      // No Error handling yet
	var table = (&tables.Downlaod{})
	table.New("https://codeload.github.com/torvalds/linux/zip/refs/heads/master", nil, nil)
	println()
	println()
	logrus.Infof("downlaod : \n %++v", table)
	println()
	println()

	table.Init()
}

