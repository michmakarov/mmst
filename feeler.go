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
	blockedRA   []string //220420 21:02
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
func (f *feeler) WriteFLog(r *http.Request, accName []byte) {
	var s string
	//if _, err = f.log.WriteString(s); err != nil {
	//	panic(fmt.Sprintf("writing into feeler log file err=%s", err.Error()))
	//}
	s = fmt.Sprintf("DATE=%s--NUM=%d--ACC=%v--URI=%s--RA=%s\n", time.Now().Format("20060102_150405"), f.getReqCount(), accName, r.RequestURI, r.RemoteAddr)

	f.mtx.Lock()
	defer f.mtx.Unlock()

	if isDebug(serverMode) {
		fmt.Print(s)
	}

	writeAllToFile(f.log, []byte(s))

}
func (f *feeler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var accRes byte
	var accName []byte
	//var aN string
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
			WriteToCommonLog(msg, -1)
		}

	}

	accName, accRes = convertCookieToAccName(r) //getCookieVal(r)

	f.WriteFLog(r, accName) //220408 08:29 (220322-account : confirmation) The feeler log fixes all incoming requests.

	if reqAnswered(w, r, accName, accRes) {
		return
	}

	if perfList.inPerforming(r) {
		s = fmt.Sprintf("%s (from %s) is in performing. You must wait", r.RequestURI, r.RemoteAddr)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(400)
		w.Write([]byte(s))
		fmt.Println(s)
		return

	}

	if r.URL.Path == "/registerme" { //220329
		if accRes == 0 { // a repeated request for registration
			//prolongeAccount(accName)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(200)
			w.Write([]byte(youHasAccAlready))
			return
		} else {
			if IPHasAccount(IPFromRA(r.RemoteAddr)) == nil {
				accName = setCookie(w)
				regAccount(accName, r)
				//f.BlockRA(r.RemoteAddr)
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(200)
				w.Write([]byte(CookieIs)) //220420 19 58 Next to cookie will be sent and receiving atempts reristration from that ip wll must be refusing until receiving a valid cookie
				return
			} else {
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(200)
				w.Write([]byte(yourIPHasAcc))
				return
			}
		}
	} else { // r.URL.Path != "/registerme"
		if accRes == 0 { //220330 15:56 there is an account; all Ok
			//f.UnblockRA(r.RemoteAddr)
			prolongeAccount(accName)
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
		//if accRes != 0 {//220422 06:51
		//	r = r.WithContext(context.WithValue(r.Context(), AccNameCtxKey, "?"))
		//} else {
		r = r.WithContext(context.WithValue(r.Context(), AccNameCtxKey, accName))
		//}
		r = r.WithContext(context.WithValue(r.Context(), URLCtxKey, r.RequestURI)) //190408
		ctx, _ = context.WithCancel(r.Context())
		r = r.WithContext(ctx)
	}

	perfList.Reg(r)
	f.h.ServeHTTP(w, r)
	perfList.Done(r)
}

/* 220416:19
func (f feeler) BlockRA(RA string) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	f.blockedRA = append(f.blockedRA, RA)
	WriteToCommonLog(fmt.Sprintf("Was blocked RA=%s", RA))
}

//220425 06:34
//What is needed to be successful in programing in first turn?
//What does it take to be successful in programing in first place?
//What does it take before all not to do foolishnes in programing ?
//What is needed, first of all, in order not to do stupid things in programming?
//First of all,I think, it is to write the purpose (and results) in natural language before the programing. So:
//It removes a RA if has found
//it panics if has not found at all or has found more than one item
func (f feeler) UnblockRA(RA string) {
	var newBlockedRA = make([]string, 0)
	var count int
	f.mtx.Lock()
	defer f.mtx.Unlock()
	//f.blockedRA = append(f.blockedRA, RA)
	for _, val := range f.blockedRA {
		if val == RA {
			count++
		} else {
			newBlockedRA = append(newBlockedRA, val)
		}
	}
	if count == 1 {
		f.blockedRA = newBlockedRA
		WriteToCommonLog(fmt.Sprintf("Was UnBblocked RA=%s", RA))
	}
	if count == 0 {
		panic(fmt.Sprintf("RA==%s was not found among blocked ones", RA))
	}

	if count > 1 {
		panic(fmt.Sprintf("RA==%s was found more than 1 blocked ones", RA))
	}
}

func (f feeler) isBlocked(RA string) bool {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	for _, val := range f.blockedRA {
		if val == RA {
			return true
		}
	}
	return false
}

//220425 12:31
func (f feeler) getBlocked() []string {
	var cpBlocked []string
	f.mtx.Lock()
	defer f.mtx.Unlock()
	for _, val := range f.blockedRA {
		cpBlocked = append(cpBlocked, val)
	}
	return f.blockedRA
}
------------*/

//220504 19:220
func (f feeler) getReqCount() int64 {
	return atomic.LoadInt64(&f.feelerCount)
}
