package tables

import (
	_ "github.com/mattn/go-sqlite3"
)

type Downlaod struct {
	IdDownlaod        int    `json:"id-downlaod"`
	IdDownlaodType   int    `json:"id-downlaod-type"`
	IdDownlaodStatus int    `json:"id-downlaod-status"`
	IdFileType       int    `json:"id-file-type"`
	WorkingFileType  string `json:"working-file-type"`
	OutputFilePath   string `json:"output-file-path"`
	Remote             string `json:"remote"`
}

 


	



