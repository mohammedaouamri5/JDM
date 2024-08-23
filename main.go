package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	api "github.com/mohammedaouamri5/JDM-back/API"
	"github.com/mohammedaouamri5/JDM-back/utile"
	log "github.com/sirupsen/logrus"

)

func main() {

	if err := utile.Init() ; err != nil {
		log.Errorln(err.Error())
	}

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

	r.POST("/download", api.Downlowd)
	r.GET("/list", api.List)
	r.GET("/setting", api.GETInfo)
	r.POST("/setting", api.POSTPath)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
