// 220107 05:33 feeler
package main

import (
	"fmt"
	"net/http"
	"os"

	//"runtime/debug"
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
	var s string
	defer func() {
		var rec interface{}
		if rec = recover(); rec != nil {
			var msg string
			msg = fmt.Sprintf("feeler panic %v (%v)", rec, getRequestBrief(r))
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			w.Write([]byte(msg))
			time.Sleep(time.Second)
			if isDebug(serverMode) {
				//fmt.Println(string(debug.Stack()))
				fmt.Printf("++++++++++++++++++++++++%s", msg)
			}
		}
	}()
	atomic.AddInt64(&f.feelerCount, 1)

	{ //220330 12:03 deleting expired eccounts
		var accs = delExpiredAccounts()
		if len(accs) != 0 {
			var msg = fmt.Sprintf("deleting expired %v", accs)
			WriteToLog(msg)
		}

	}

	accName, accRes = getAccount2(r) //getCookieVal(r)

	logMess = fmt.Sprintf("%s:%d--ACC=%s--URI=%s--RA=%s\n", time.Now().Format("20060102_150405"), f.feelerCount, accName, r.RequestURI, r.RemoteAddr)
	f.WriteFLog(logMess) // (220322-account : confirmation) The feeler log fixes all incoming requests.
	if isDebug(serverMode) {
		fmt.Println(logMess)
	}

	if r.URL.Path == "/q" {
		s = fmt.Sprintf("There is /q debug request; accName=%s; accres=%d", accName, accRes)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(s))
		fmt.Println(s)
		return
	}

	if r.URL.Path == "/registerme" { //220329
		if accRes == 0 { // a repeated request for registration
			prolongeAccount(accName)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(200)
			w.Write([]byte(CookieIsStill))
			return
		} else {
			accName = setCookie(w)
			regAccount(accName, r)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(200)
			w.Write([]byte(CookieIs))
			return
		}
	} else {
		if accRes == 0 { //220330 15:56 there is an account; all Ok
			goto toMultiplexer
		} else { //the 400  will be passed to the client
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(400)
			w.Write([]byte(noCookieMess))
			return
		}
	}

toMultiplexer:
	f.h.ServeHTTP(w, r)
}
