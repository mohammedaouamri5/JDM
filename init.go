package main

import (
	"io"
	"time"

	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func InitLog() {
	// Set custom formatter
	// Set output to both stdout and the log file
	multiWriter := io.MultiWriter(colorable.NewColorableStdout())
	log.SetOutput(multiWriter)

	field_map := log.FieldMap{
		// log.FieldKeyTime:  "@timestamp",
		log.FieldKeyMsg:  "@message",
		log.FieldKeyFunc: "@caller",
	}

	// Set the formatter to include timestamp and caller information
	textformatter := log.TextFormatter{
		FullTimestamp:          true,
		ForceColors:            true,
		ForceQuote:             true,
		PadLevelText:           true,
		DisableLevelTruncation: false,
		FieldMap:               field_map,
		TimestampFormat:        time.DateTime,
	}

	log.SetFormatter(&textformatter)

	// Enable reporting caller information
	log.SetReportCaller(true)

}
