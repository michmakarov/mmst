// log
//220308 08:19
package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var lf *os.File
var logger *log.Logger

func init() {
	var err error
	var pref = time.Now().Format(timeFormat)
	if lf, err = os.Create(logFileName); err != nil {
		panic(fmt.Sprintf("Error of creating a log file=%s\n", err.Error()))
	}
	logger = log.New(lf, pref, log.Lshortfile)
}
func WriteToLog(msg string) {
	logger.Output(2, msg)
}
