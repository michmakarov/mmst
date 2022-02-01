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

	//"html/template"
	"sync/atomic"
	"syscall"
	"time"
	//_ "github.com/lib/pq"
)

const timeFormat = "20060102_15:04:05"
const maxLetters = 10
const maxChars = 3000

var appName = "---mmsite from 220131_1448---" //It is assigned in actual value by b.sh

var RequestCounter uint32 //increased atomically by each handler

var exitByError = true

//220128 10:32
//The last digit (most low) defines that there is debugging (0 - not, 1 - is)
//second - protocol (0 - http, 1 - https); so default 0 says http without debugging
// see setArgs(), isDebug(int)bool, and isHTTPS(int)bool in utils.go
// third - sending sms when receiving a letter; see isDebug(mode int) bool, isHTTPS, isSms
//So 1: http without sms with debug
// 101 : http with sms with debug
var serverMode int

var flr *feeler //220108 08:22

var lang string = "ru"
var passWord string = "***not assigned yet***"

//Those are default values. See func setArg()
//var PG_CONN_STR = "postgres://kot_user:1qazXSW@@localhost:5433/mak_docker"
var listeningAddr = "0.0.0.0:8080"

func main() {

	setArgs()

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

	mx.HandleFunc("/e_2", e_2Handler)

	flr = createFeeler(mx)
	srv := &http.Server{
		Handler:      flr,
		Addr:         listeningAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	srv.RegisterOnShutdown(onShutDown)

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
	fmt.Printf("%v\n", appName)
	fmt.Printf("The server will be listen at %v, mode=%v\n\n", listeningAddr, serverMode)
	//-----220106 18:11

	// Start Server
	go func() {
		switch isHTTPS(serverMode) {
		case false:
			if err := srv.ListenAndServe(); err != nil {
				fmt.Printf("The server refused to work with error:%v\n", err.Error())
			}
		case true:
			if err := srv.ListenAndServeTLS("cert", "key"); err != nil {
				fmt.Printf("The server refused to work with error:%v\n", err.Error())
			}
		} // switch
	}()

	// Graceful Shutdown
	waitForShutdown(srv)

}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

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
