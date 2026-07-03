package logger

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
)

func Init(name, level string) {
	log.SetPrefix(name + "\t")

	Info = log.New(
		os.Stdout,
		name+" [INFO] ",
		log.Ldate|log.Ltime,
	)

	Error = log.New(
		os.Stderr,
		name+" [ERROR] ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)

	var debugWriter io.Writer = io.Discard

	if level == "debug" {
		debugWriter = os.Stdout
	}

	Debug = log.New(
		debugWriter,
		name+" [DEBUG] ",
		log.Ldate|log.Ltime,
	)

}
