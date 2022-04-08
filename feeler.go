// 220107 05:33 feeler
package main

import (
	"context"
	"fmt"

	//"io"
	"net/http"
	"os"

	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

//220331 08:39 CtxParType is the type of context parameters which assigning to incoming requests
type ctxParType string

const (
	//Next keys are established by func (f *feeler) ServeHTTP
	AccNameCtxKey ctxParType = "AccName"
	ReqNumCtxKey  ctxParType = "RecNum"
	URLCtxKey     ctxParType = "URL"
)

type feeler struct {
	h           http.Handler // invoking handler
	feelerCount int64
	log         *os.File
	logFileName string
	mtx         *sync.Mutex
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
	f.logFileName = fLogFileName
	if _, err = f.log.WriteString(fmt.Sprintf("Start %s\n", time.Now().Format("20060102_150405"))); err != nil {
		fmt.Printf("inserting start record into feeler log file err=%s", err.Error())
		os.Exit(1)
	}
	f.mtx = &sync.Mutex{}

	return f
}

//220325 11:06
//func (f *feeler) WriteFLog(r *http.Request) {
//220408 08:29
func (f *feeler) WriteFLog(r *http.Request, aN string) {
	var s string
	//if _, err = f.log.WriteString(s); err != nil {
	//	panic(fmt.Sprintf("writing into feeler log file err=%s", err.Error()))
	//}
	s = fmt.Sprintf("DATE=%s--NUM=%d--ACC=%s--URI=%s--RA=%s\n", time.Now().Format("20060102_150405"), f.feelerCount, aN, r.RequestURI, r.RemoteAddr)

	f.mtx.Lock()
	defer f.mtx.Unlock()

	if isDebug(serverMode) {
		fmt.Print(s)
	}

	writeAllToFile(f.log, []byte(s))

}
func (f *feeler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var accRes byte
	var accName, aN string
	//var logMess string
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
				fmt.Printf("++++++++++++++++++++++++%s\n", msg)
				fmt.Println(string(debug.Stack()))
				fmt.Printf("++++++++++++++++++++++++\n")
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

	if accRes == 0 {
		aN = fmt.Sprintf("%v", []byte(accName))
	} else {
		aN = fmt.Sprintf("accRes==%d", accRes)
	}

	f.WriteFLog(r, aN) //220408 08:29 (220322-account : confirmation) The feeler log fixes all incoming requests.
	//if isDebug(serverMode) {
	//	fmt.Println(logMess)
	//}

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
			regAccount([]byte(accName), r)
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
	{ //220331 09:19
		var ctx context.Context

		r = r.WithContext(context.WithValue(r.Context(), ReqNumCtxKey, strconv.FormatInt(f.feelerCount, 10)))
		if accRes != 0 {
			r = r.WithContext(context.WithValue(r.Context(), AccNameCtxKey, "?"))
		} else {
			r = r.WithContext(context.WithValue(r.Context(), AccNameCtxKey, accName))
		}
		r = r.WithContext(context.WithValue(r.Context(), URLCtxKey, r.RequestURI)) //190408
		ctx, _ = context.WithCancel(r.Context())
		r = r.WithContext(ctx)
	}

	f.h.ServeHTTP(w, r)
}
