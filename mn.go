//220104 10:43 See ../Readme.txt
//it is thought as proving of skills for a potential employer
//Prototype Progects/docker/variant1/http_srv_211025
package main

import (
	"context"
	"fmt"

	"net/http"

	"os"
	"os/signal"

	"sync/atomic"
	"syscall"
	"time"

	"net"
	"runtime"
)

//220426 14:16 See also the history and (this date record) SaveConnInContext, GetConn, ConnContextKey
type contextKey struct {
	key string
}

//240427 08:51 See saveConnChanging
type CurrentConnection struct {
	conn  net.Conn
	state http.ConnState
}

var ConnContextKey = &contextKey{"mmst-conn"}

const timeFormat = "20060102_150405"
const serverMaxMode = 99999
const maxLetters = 10
const maxChars = 3000

var currConn = &CurrentConnection{}

//var appName = "---mmsite from 220131_1448---" //It is assigned in actual value by b.sh

var versionInfo = "version not assigned yet" //220202 06:26

var RequestCounter uint32 //increased atomically by each handler

var exitByError = true

//220128 10:32
//The last digit (most low) defines that there is debugging (0 - not, 1 - is)
//second - protocol (0 - http, 1 - https); so default 0 says http without debugging
// see setArgs(), isDebug(int)bool, and isHTTPS(int)bool in utils.go
// third - sending sms when receiving a letter; see isDebug(mode int) bool, isHTTPS, isSms
//So 1: http without sms with debug
// 101 : http with sms with debug
//220421 14:58 four - hiding the output to a file
//that is 1001 - out to file, without sms, http, debug; see also func isHidingOut(mode int) bool;func digInPos(mode int, pos int) int
var serverMode int

var flr *feeler                  //220108 08:22
var maxFrontLogSize = 1000000000 //220408 17:38 100Mb

//var lang string = "ru"
var passWord string = "***not_assigned_yet***" //220405 10:41 220420 10:34:for there is obscurity with sending it as parameter

//Those are default values. See func setArg()
//var PG_CONN_STR = "postgres://kot_user:1qazXSW@@localhost:5433/mak_docker"
var listeningAddr = "0.0.0.0:8080"

//220302 17:19
//220308 08:53
var accountsFileName string = "accountList.txt"
var accountsFileMaxSize = 10000 //bytes

//220308 08:23
var generalLogFileName string = "general.log"

//220311 15:00
var httpLogFileName string = "http.log"

func main() {

	if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
		fmt.Printf("This program intends to work only under linux/arm64, but here is %s and %s\n", runtime.GOOS, runtime.GOARCH)
		return
	}

	setArgs()
	redirectStd()
	removeOldLogs()
	creteGeneralHttpLogs()

	mx := http.DefaultServeMux

	mx.HandleFunc("/", indHandler)
	mx.HandleFunc("/ind", indHandler)
	mx.HandleFunc("/favicon.ico", faviconHandler)
	mx.HandleFunc("/showMemStat", showMemStatHandler)
	mx.HandleFunc("/stop", stopHandler)
	mx.HandleFunc("/about", aboutHandler)
	mx.HandleFunc("/changeLang", changeLangHandler)
	mx.HandleFunc("/mmsite_script", mmsite_scriptHandler)
	mx.HandleFunc("/aboutAuthor", aboutAuthorHandler)
	mx.HandleFunc("/myFamily", myFamilyHandler)
	mx.HandleFunc("/history", historyHandler)
	mx.HandleFunc("/accounts", accountsHandler)
	mx.HandleFunc("/sms", smsHandler)
	mx.HandleFunc("/letter", letterHandler)
	mx.HandleFunc("/css", cssHandler)

	mx.HandleFunc("/myAccount", myAccountHandler)
	mx.HandleFunc("/showFeelerLog", showFeelerLogHandler)
	mx.HandleFunc("/help", helpHandler)
	mx.HandleFunc("/showGeneralLog", showGeneralLogHandler)
	mx.HandleFunc("/longOper", longOperHandler)
	mx.HandleFunc("/main", mainHandler)

	mx.HandleFunc("/e_2", e_2Handler)

	flr = createFeeler(mx)
	srv := &http.Server{
		Handler:      flr,
		Addr:         listeningAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     servLogger,
		ConnContext:  SaveConnInContext, //220426
		ConnState:    saveConnChanging,  //220427 08:41
	}

	srv.RegisterOnShutdown(onShutDown)
	/* 220202 08:22 All this now will do the b.sh
	//220106 18:11
	var err error
	var proc *os.Process
	var set_ind_args []string
	var set_ind_attr os.ProcAttr = os.ProcAttr{Files: []*os.File{nil, os.Stdout, os.Stderr}}
	var procSt *os.ProcessState

	if proc, err = os.StartProcess("./set_ind.sh", set_ind_args, &set_ind_attr); err != nil {
		fmt.Printf("Starting ./set_ind.sh err\n%v\n", err.Error())
		os.Exit(1)
	}
	if procSt, err = proc.Wait(); err != nil {
		fmt.Printf("Waiting of ./set_ind.sh err\n%v\n", err.Error())
		os.Exit(1)
	} else {
		fmt.Printf("The ind.html was successly changed; procSt.String():\n%v\n", procSt.String())
	}
	//time.Sleep(time.Second)
	*/
	fmt.Printf("ver: %v\n", versionInfo)
	//-----220106 18:11

	// Start Server
	go func() {
		switch isHTTPS(serverMode) {
		case false:
			fmt.Printf("The server (http) will be listen at %v, mode=%v\n\n", listeningAddr, serverMode)
			if err := srv.ListenAndServe(); err != nil {
				fmt.Printf("The server refused to work with error:%v\n", err.Error())
			}
		case true:
			fmt.Printf("The server (https) will be listen at %v, mode=%v\n\n", listeningAddr, serverMode)
			if err := srv.ListenAndServeTLS("cert", "key"); err != nil {
				fmt.Printf("The server refused to work with error:%v\n", err.Error())
			}
		} // switch
	}()

	// Graceful Shutdown
	waitForShutdown(srv)

	//220302 16:58
	//saveAccounts("accoutList.txt")

}

func waitForShutdown(srv *http.Server) {
	var incomeSig os.Signal
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	incomeSig = <-interruptChan
	WriteToCommonLog(fmt.Sprintf("waitForShutdown: a signal was received: %s", incomeSig.String()), -1)

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)

	fmt.Printf("Shutting down by user signal\n")
	os.Exit(0)
}

func onShutDown() {
	fmt.Printf("onShutdown starts ...\n")
	fmt.Printf(" ... onShutdown ends \n")
}

func getRequestCounter() uint32 {
	return atomic.AddUint32(&RequestCounter, 1)
}

//220414 05:56
func redirectStd() {
	if !(isHidingOut(serverMode)) {
		return
	}
	var err error
	var nStd *os.File
	//return //220420 18:31
	if nStd, err = os.Create("out.txt"); err != nil {
		panic(fmt.Sprintf("redirectStd:creating out.txt err=%s", err.Error()))
	}
	os.Stdout = nStd
	os.Stderr = nStd
}

func SaveConnInContext(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, ConnContextKey, c)
}

func GetConn(r *http.Request) net.Conn {
	return r.Context().Value(ConnContextKey).(net.Conn)
}

//220427 08:30
func saveConnChanging(conn net.Conn, cs http.ConnState) {
	var msg string
	if currConn.conn != conn {
		msg = fmt.Sprintf("There is changing state of connection from %s (new state=%s)", conn.RemoteAddr().String(), cs.String())
		WriteToCommonLog(msg, -1)
	}
	currConn.conn = conn
	currConn.state = cs
}
