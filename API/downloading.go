package api

import (
	"net/http"

	"github.com/gin-gonic/gin" 
	"github.com/mohammedaouamri5/JDM-back/downloader"
	vector "github.com/mohammedaouamri5/vector/vector"
	log "github.com/sirupsen/logrus"
)

var Db = vector.New[*downloader.FILE](5, 0.5, [](*downloader.FILE){});



 

func Download(c *gin.Context) {
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
 
	Db.Push(&file)
	
	
	
	if err := file.MkeConfig(5); err != nil {
		log.Errorln(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
	"ID": file.Id.String(),
	}); 
}
func List(c *gin.Context) {

	for i := 0; i < Db.Size(); i++ {
		(*(*Db.Data)[i]).ReadFromMetaData(downloader.Cfgjson((*(*Db.Data)[i]))) 
	}
	// log.Infof("\n\tList : %+v\n\t " , Db.Data)
	c.JSON(http.StatusOK, gin.H{"data": Db.Data })
}

func Pause_Unpause(c *gin.Context)  {
	type RequestPOSTPath_Unpause struct {
		ID string `json:"ID" binding:"required"`	
	}
	var req RequestPOSTPath_Unpause
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	for i := 0; i < Db.Size(); i++ {
		if (*(*Db.Data)[i]).Id.String() == req.ID {
			log.Infof("\n\tList : %+v\n\t " , (*(*Db.Data)[i]))
			(*(*Db.Data)[i]).IsPause = !(*(*Db.Data)[i]).IsPause
			log.Infof("\n\tList : %+v\n\t " , (*(*Db.Data)[i]))
			c.JSON(http.StatusOK , gin.H{
				"message": "success",
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "not exist"})	

	
}




    










