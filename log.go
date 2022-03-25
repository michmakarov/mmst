// log
//220308 08:19
package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var lf, sLf *os.File
var logger, servLogger *log.Logger

func init() {
	var err error
	var pref = time.Now().Format(timeFormat)
	if lf, err = os.Create(logFileName); err != nil {
		panic(fmt.Sprintf("Error of creating a log file=%s\n", err.Error()))
	}
	if sLf, err = os.Create(httpLogFileName); err != nil {
		panic(fmt.Sprintf("Error of creating a HTTP (SERVER) log file=%s\n", err.Error()))
	}
	logger = log.New(lf, pref, log.Lshortfile)
	servLogger = log.New(sLf, pref, log.Lshortfile)
}
func WriteToLog(msg string) {
	logger.Output(2, msg)
}
func WriteToServLog(msg string) {
	servLogger.Output(2, msg)
}
