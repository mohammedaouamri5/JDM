package main

import (
	_ "github.com/mattn/go-sqlite3"
	HTTP "github.com/mohammedaouamri5/JDM-back/Downlaod/HTTP"
	. "github.com/mohammedaouamri5/JDM-back/db"
	"github.com/mohammedaouamri5/JDM-back/tables"
)

func main() {
	InitLog() // No Error handling yet
	DB()      // No Error handling yet
	tables.State{}.Pull()
	var table = (&tables.Downlaod{})
	table.New("http://212.183.159.230/5MB.zip", nil, nil)
	println()
	println()
	println()
	println()
	table.Init()
	HTTP.Downlaod(*table)
}
