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
func (f *feeler) WriteFLof(r *http.Request) {
	var err error
	if _, err = f.log.WriteString(fmt.Sprintf("%s:%d--URI=%s--RA=%s\n", time.Now().Format("20060102_150405"), f.feelerCount, r.RequestURI, r.RemoteAddr)); err != nil {
		panic(fmt.Sprintf("writing into feeler log file err=%s", err.Error()))
	}
}
func (f *feeler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	f.WriteFLof(r)
	regAccount(r)

	if isDebug(serverMode) {
		fmt.Printf("---------feeler: RA=%s; URI=%s; lang=%s\n", r.RemoteAddr, r.RequestURI, getLang(r))
	}

	f.h.ServeHTTP(w, r)
}
