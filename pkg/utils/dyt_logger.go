package utils

import (
	"log"
	"os"
)

// func ErrorHandler(err error, message string) error {
// 	// Code ANSI pour texte rouge: \033[31m et reset: \033[0m
// 	errorLogger := log.New(os.Stderr, "\033[31mERROR\033[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
// 	errorLogger.Println(message, err)
// 	return fmt.Errorf("%s", message)
// }

func PrintHanlder(message string) {
	printLogger := log.New(os.Stdout, "\033[32mINFO\033[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
	printLogger.Println(message)
}
