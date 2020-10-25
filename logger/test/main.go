package main

import (
	"github.com/hanjingo/logger"
)

//for win: go build -o log.exe main.go
func main() {
	log := logger.GetDefaultLogger()
	logger.DefaultConsoleWriter.SetValid(false)
	for i := 0; i < 100000; i++ {
		log.Info("info")
		log.Debug("debug")
		log.Fatal("fatal")
		log.Notice("notice")
		log.Error("error")
	}
}
