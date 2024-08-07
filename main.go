package main

import ( 

	"github.com/mohammedaouamri5/JDM-back/downloader"
	"github.com/mohammedaouamri5/JDM-back/utile"
	log "github.com/sirupsen/logrus"
)
 


func main()   { 
	if err := utile.Init() ; err != nil {
		log.Errorln(err.Error())
	 	println("\n\nEXIT 1") 	   
		return 
	}

	path := string("./tmp")
 
	var file downloader.FILE    
	file.Constructor(
		"https://archive.archlinux.org/packages/d/dolphin/dolphin-24.07.80-1-x86_64.pkg.tar.zst" , 
		"dolphin-24.07.80-1-x86_64.pkg.tar.zst" , &path )  
		log.Print(file.Download(10) )

	println("\n\nEXIT 0") 	      
}
  