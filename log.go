// log
//220308 08:19
package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var logger, servLogger *log.Logger

func creteGeneralHttpLogs() {
	var lf, sLf *os.File
	var err error
	var pref = time.Now().Format(timeFormat) + "---"
	if lf, err = os.Create(generalLogFileName); err != nil {
		panic(fmt.Sprintf("Error of creating a log file=%s\n", err.Error()))
	}
	if sLf, err = os.Create(httpLogFileName); err != nil {
		panic(fmt.Sprintf("Error of creating a HTTP (SERVER) log file=%s\n", err.Error()))
	}
	logger = log.New(lf, pref, log.Lshortfile)
	servLogger = log.New(sLf, pref, log.Lshortfile)
	WriteToCommonLog("CommonLog Started\n", -1)
	WriteToServLog("ServLog Started\n")
}
func WriteToCommonLog(msg string, calldepth int) {
	var err error
	if calldepth < 0 || calldepth > 4 {
		calldepth = 2
	}
	if err = logger.Output(calldepth, msg); err != nil {
		panic(fmt.Sprintf("WriteToCommonLog err=%s\n", err.Error()))
	}
}
func WriteToServLog(msg string) {
	servLogger.Output(2, msg)
}
