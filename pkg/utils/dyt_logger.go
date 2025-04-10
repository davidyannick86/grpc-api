package utils

import (
	"log"
	"os"
)

func PrintHanlder(message string) {
	printLogger := log.New(os.Stdout, "\033[32mINFO\033[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
	printLogger.Println(message)
}
