package utile

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)




func initLogger() error {


	logrus.SetOutput(colorable.NewColorableStdout())

	// Set the formatter to include timestamp and caller information
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	// Enable reporting caller information
	logrus.SetReportCaller(true) 
	return nil 

}

func Init() error {
	

	if err := initLogger() ; err != nil {
		logrus.Error(err.Error())
		return err 
	}

	return nil

}