package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	//"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	api "github.com/mohammedaouamri5/JDM-back/API"
	"github.com/mohammedaouamri5/JDM-back/downloader"
	"github.com/mohammedaouamri5/JDM-back/utile"
	log "github.com/sirupsen/logrus"
)

func Read_db() {
	dir := "./data"

	// Open the directory
	d, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()

	// Read the directory contents
	files, err := d.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			// Remove the extension
			nameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			IsExist := false
			for _, value := range *api.Db.Data {
				if (*value).Id.String() == nameWithoutExt {
					IsExist = true
					print("\n\nexist\n")
				}
			}
			if !IsExist {
				print("\n\nnot exist\n")
				file := downloader.FILE{Id: uuid.MustParse(nameWithoutExt)}
				file.ReadFromMetaData(downloader.Cfgjson(file))
				api.Db.Push(&file)
			}
		}
	}
}

func main() {

	if err := utile.Init(); err != nil {
		log.Errorln(err.Error())
	}
	go func() {
		for { 
			Read_db()
			log.Info(api.Db)
			log.Info(api.Db.Size())
			if api.Db.Size() > 0 {
				for i := 0; ; i = (i + 1) % (api.Db.Size()) {
					if (*(*api.Db.Data)[i]).IsDone == false && (*(*api.Db.Data)[i]).IsPause == false {
						(*(*api.Db.Data)[i]).Download(15)
					}
				}
			}
			time.Sleep(15 * time.Second)
		}
	}() 

	log.Info("test looger")

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Routing
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/download", api.Download)
	r.GET("/list", api.List)
	r.GET("/setting", api.GETInfo)
	r.POST("/setting", api.POSTPath)
	r.POST("/pause", api.Pause_Unpause)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
