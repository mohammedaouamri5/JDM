package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohammedaouamri5/JDM-back/downloader"
	vector "github.com/mohammedaouamri5/vector/vector"
	log "github.com/sirupsen/logrus"
)

var vec = vector.New[downloader.FILE](5, 0.5, []downloader.FILE{});





func Downlowd(c *gin.Context) {
	type Download_file struct {
		Url  string `json:"url" binding:"required"`
		Name  string `json:"name" binding:"required"`
		Path string `json:"path" `
	}
	var file_tmp Download_file
	var file downloader.FILE
	
	if err := c.ShouldBindJSON(&file_tmp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if  file_tmp.Path == "" {
		file.Constructor(file_tmp.Url, file_tmp.Name, nil) 
	} else {
		file.Constructor(file_tmp.Url, file_tmp.Name, &file_tmp.Path) 
	}
 
	vec.Push(file)
	
	
	
    start := time.Now()
	if err := file.Download(5); err != nil {
		log.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	elapsed := time.Since(start)
	hours := int(elapsed.Hours())
	minutes := int(elapsed.Minutes()) % 60
	seconds := int(elapsed.Seconds()) % 60
	
	c.JSON(http.StatusOK, gin.H{
		"took": gin.H{
			"HH": fmt.Sprintf("%02d", hours),
			"MM": fmt.Sprintf("%02d", minutes),
			"SS": fmt.Sprintf("%02d", seconds),
		},
		"elapsed": elapsed,
	}) 
}
func List(c *gin.Context) {

	for i := 0; i < vec.Size(); i++ {
		(*vec.Data)[i].ReadFromMetaData(downloader.Cfgjson((*vec.Data)[i])) 
	}
	// log.Infof("\n\tList : %+v\n\t " , vec.Data)
	c.JSON(http.StatusOK, gin.H{"data": vec.Data })

}







    










