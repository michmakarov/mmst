// 220107 05:33 feeler
package main

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

type feeler struct {
	h           http.Handler // invoking handler
	feelerCount int64
	log         *os.File
}

func createFeeler(h http.Handler) (f *feeler) {
	var err error
	f = &feeler{}
	f.h = h
	var fLogFileName = "Feeler" + time.Now().Format("20060102_150405") + ".log"
	if f.log, err = os.Create(fLogFileName); err != nil {
		fmt.Printf("Creating feeler log file err=%s", err.Error())
		os.Exit(1)
	}
	if _, err = f.log.WriteString(fmt.Sprintf("Start %s\n", time.Now().Format("20060102_150405"))); err != nil {
		fmt.Printf("inserting start record into feeler log file err=%s", err.Error())
		os.Exit(1)
	}

	return f
}

//220325 11:06
//func (f *feeler) WriteFLog(r *http.Request) {
func (f *feeler) WriteFLog(s string) {
	var err error

	//if _, err = f.log.WriteString(fmt.Sprintf("%s:%d--URI=%s--RA=%s\n", time.Now().Format("20060102_150405"), f.feelerCount, r.RequestURI, r.RemoteAddr)); err != nil {
	//	panic(fmt.Sprintf("writing into feeler log file err=%s", err.Error()))
	//}
	if _, err = f.log.WriteString(s); err != nil {
		panic(fmt.Sprintf("writing into feeler log file err=%s", err.Error()))
	}
}
func (f *feeler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var accRes byte
	var accName string
	var logMess string
	defer func() {
		var rec interface{}
		if rec = recover(); rec != nil {
			var msg string
			msg = fmt.Sprintf("feeler panic %v (%v)", rec, getRequestBrief(r))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			w.Write([]byte(msg))
		}
	}()
	atomic.AddInt64(&f.feelerCount, 1)
	accName, accRes = getCookieVal(r)

	logMess = fmt.Sprintf("%s:%d--URI=%s--RA=%s\n", time.Now().Format("20060102_150405"), accName, f.feelerCount, r.RequestURI, r.RemoteAddr)
	f.WriteFLog(logMess) // (220322-account : confirmation) The feeler log fixes all incoming requests.
	if isDebug(serverMode) {
		fmt.Printf("--feeler: RA=%s; URI=%s; account=%s\n", r.RemoteAddr, r.RequestURI, accName)
	}

	if r.URL.Path == "/" || r.URL.Path == "/ind" {
		if accRes == 0 { // a repeated request of the index page
			goto toMultiplexer
		} else { //There is possibility that there is an appropriate account. It must be deleted with creating a new one
			delAccount(accName) // be removed if exists
			accName = setCookie(w)
			regAccount(accName, r)
			WriteToLog(fmt.Sprintf("%s was removed (maybe, err %d) and %s was registered", accName, accRes, accountName(r)))
			goto toMultiplexer
		}
	}
	//regAccount(r)

toMultiplexer:
	f.h.ServeHTTP(w, r)
}
