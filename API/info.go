package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/mohammedaouamri5/JDM-back/utile"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func GETInfo(c *gin.Context) {
	// Convert InfoS to JSON
	infoBytes, err := json.Marshal(utile.InfoS)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process data"})
		return
	}

	// Convert JSON to map
	var infoMap map[string]interface{}
	if err := json.Unmarshal(infoBytes, &infoMap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process data"})
		return
	}

	// Return the map as JSON
	c.JSON(http.StatusOK, infoMap)
}

func POSTPath(c *gin.Context) {
	type RequestPOSTPath struct {
		PATH string `json:"PATH" binding:"required"`
	}
	var requestPOSTpath RequestPOSTPath

	if err := c.ShouldBindJSON(&requestPOSTpath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save PATH
	utile.InfoS.PATH = requestPOSTpath.PATH

	if err := utile.SaveInfo(); err != nil {
		log.Errorln(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
