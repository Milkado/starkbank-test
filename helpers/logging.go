package helpers

import (
	"log"
	"os"
	"path/filepath"
)

func LogError(logFile string, message string) {
	absPath, err := filepath.Abs(logFile)
	if err != nil {
		log.Fatal("error openning the file: " + err.Error())
	}

	file, err := os.OpenFile(absPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("error openning the file: " + err.Error())
	}

	log.SetOutput(file)

	log.Println(message)
}
