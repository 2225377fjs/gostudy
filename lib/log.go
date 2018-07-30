package lib

import (
	"log"
	"os"
	"time"
)

// 实现异步的文件log

var logger *log.Logger
var contents chan []interface{}


func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func init(){
	fileName := "./logs/" + time.Now().Format("2006-01-02") + ".log"
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	logger = log.New(file, "", log.LstdFlags|log.Llongfile)
	contents = make(chan []interface{}, 500)
	go doLog()
}


func doLog() {
	for {
		select {
		case content := <- contents:
			logger.Print(content ...)
			log.Print(content ...)
		}
	}
}

func Log(a ...interface{}) {
	contents <- a
}