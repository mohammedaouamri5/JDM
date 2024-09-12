package api

import (
	"net/http"

	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	downlaod "github.com/mohammedaouamri5/JDM-back/Downlaod"
	"github.com/mohammedaouamri5/JDM-back/tables"
	"github.com/sirupsen/logrus"
)


type POSTDownlaodHTTPInterface struct {
    Remote string `json:"remote" binding:"required"`
    Name   string `json:"name"`
    Outdir string `json:"outdir"`
}
func (me *POSTDownlaodHTTPInterface) ToDownload() tables.Downlaod {
	return (&tables.Downlaod{}).New(
		(*me).Remote,
		"",
		(*me).Outdir,
		(*me).Name,
	)
}

func POSTDownlaodHTTP(gctx *gin.Context) {
	post_downlaod := POSTDownlaodHTTPInterface{}
	if err := gctx.ShouldBindJSON(&post_downlaod); err != nil {
		logrus.Error(err.Error())
		gctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"Error": err.Error(),
			},
		)
		return
	}
	logrus.Info(post_downlaod)
	downlaod := (&post_downlaod).ToDownload()
	downlaod.Init()
	gctx.JSON(
		http.StatusCreated,
		structs.Map(downlaod),
	)
}



type POSTPauseInterface struct {
    Id int `json:"id"`
}
func POSTPause(gctx *gin.Context){
	post_pause := POSTPauseInterface{}
	if err := gctx.ShouldBindJSON(&post_pause); err != nil {
		logrus.Error(err.Error())
		gctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"Error": err.Error(),
			},
		)
		return
	}

	downlaod.Pause(post_pause.Id)

}
