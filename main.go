package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	downlaod "github.com/mohammedaouamri5/JDM-back/Downlaod"
	"github.com/mohammedaouamri5/JDM-back/api"
	. "github.com/mohammedaouamri5/JDM-back/db"
	"github.com/mohammedaouamri5/JDM-back/tables"
)

func main() {
	InitLog() // No Error handling yet
	DB()      // No Error handling yet
	tables.State{}.Pull()
	go downlaod.DownloadLiciner()	
	router := gin.Default()
	router.GET("/BRUH" , func (c * gin.Context ) { c.JSON(http.StatusOK , gin.H{"messege": "BRUH"}) }) 
	router.POST("/Download/HTTP/start" , api.POSTDownlaodHTTP)
	router.POST("/Download/HTTP/pause" , api.POSTPause)
	router.Run()
}
