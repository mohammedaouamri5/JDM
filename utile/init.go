package utile

import (
	"encoding/json"
	"io"
	"io/ioutil"
	os "os"
	"time"

	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

type InfoT struct {
	PATH string `json:"PATH" binding:"required"`
}

var InfoS = InfoT{}

func initInfo() error {
	file, err := os.Open("info.json")
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Errorln("Error reading file:", err.Error())
		return err
	}
 
	// Check if the file is empty
	if len(content) == 0 {
		log.Errorln("info.json is empty")
		return nil // or return an error if an empty file is critical
	}

	err = json.Unmarshal(content, &InfoS)
	if err != nil {
		log.Errorln("Error unmarshalling JSON:", err.Error())
		return err
	}

	return nil
}


func SaveInfo() error {
 

	file, err := os.OpenFile("info.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Errorf("Error opening file: %v", err)
		return err
	}
	defer file.Close()

	content, err := json.MarshalIndent(InfoS, "", "  ")
	if err != nil {
		log.Errorln("\n\t", "Error marshalling JSON:", err.Error())
		return err
	}
	_, err = file.Write(content)
	if err != nil {
		log.Errorln("\n\t", "Error writing to file:", err.Error())
		return err
	}

	defer file.Close()
	return nil 
}  

func initLogger() error {

	


	// Set output to both stdout and the log file
	multiWriter := io.MultiWriter(colorable.NewColorableStdout())
	log.SetOutput(multiWriter)
	
	
	
	field_map := log.FieldMap{
		// log.FieldKeyTime:  "@timestamp",
		log.FieldKeyMsg:   "@message",
		log.FieldKeyFunc:  "@caller",
		
	}

	// Set the formatter to include timestamp and caller information
	textformatter := log.TextFormatter{
		FullTimestamp			:	true  ,
		ForceColors				:	true  ,
		ForceQuote				:	true  ,
		PadLevelText			:	true  , 
		DisableLevelTruncation	:	false ,
		FieldMap			    :   field_map ,
		TimestampFormat:time.DateTime,
	}
	
	log.SetFormatter(&textformatter)

	// Enable reporting caller information
	log.SetReportCaller(true)

	return nil
}

func Init() error {

	if err := initLogger(); err != nil {
		log.Error(err.Error())
		return err
	}

	if err := initInfo(); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil

}
