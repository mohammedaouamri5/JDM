package tables 

import (
	_ "github.com/mattn/go-sqlite3"
)

type Packet struct {
	IdPacket int    `json:"id-packet"`
	IdDownlaod int    `json:"id-downlaod"`
	Start int64 `json:"start"`
	End int64 `json:"end"`
	IdState int `json:"id-state"`

}


