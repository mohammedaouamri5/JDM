package tables

import (
	_ "github.com/mattn/go-sqlite3"
)



type FileType struct {
	IdFileType int    `json:"id-file-type"`
	Name       string `json:"name"`	
}




