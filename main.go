package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/mohammedaouamri5/JDM-back/API" 
	"github.com/mohammedaouamri5/JDM-back/utile"

	log "github.com/sirupsen/logrus"
)





func main() {
	if err := utile.Init(); err != nil {
		log.Errorln(err.Error())
		println("\n\nEXIT 1")
		return
	}

	

	r := gin.Default()
	/* the routring */ {

		r.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
		
		r.POST("/download", api.Downlowd)
		r.GET( "/list"    , api.List    )
		
     	}
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	return 
}
