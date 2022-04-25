// utils
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	//"os/exec"
	"path/filepath"

	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func qqqutils() {
	fmt.Println("qqqutils:Hello World!")
}

func memStatStr(start time.Time) (res string) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	res = fmt.Sprintf("Alloc=%vmgb;<br>Sys=%vmgb;<br> dur=%v;",
		ms.Alloc/1000000, ms.Sys/1000000, time.Since(start))
	return
}

//211025 - 211027_13:41
//220105 13:18
func setArgs() {
	for i := 1; i < len(os.Args); i++ {
		var splitedArg = strings.Split(os.Args[i], "=")
		if len(splitedArg) != 2 {
			fmt.Printf("It are allowed arguments with format name=value, but is %v\n", os.Args[i])
			os.Exit(1)
		}
		switch splitedArg[0] {
		case "la":
			listeningAddr = splitedArg[1]
		case "pw":
			passWord = splitedArg[1]
		case "mode":
			serverMode, _ = strconv.Atoi(splitedArg[1])
			if serverMode < 0 || serverMode > serverMaxMode {
				fmt.Printf("The mode must be > 0 and < serverMaxMode\n")
				os.Exit(1)
			}
		default:
			fmt.Printf("There is not allowed option %v\n", splitedArg[0])
			os.Exit(1)
		}
	}
}

//220107 06:46 For func (f *feeler) ServeHTTP (handling panic)
func getRequestBrief(r *http.Request) string {
	return fmt.Sprintf("%v from %v", r.RequestURI, r.RemoteAddr)
}

//220113 05:20 A analog of http.ServeFile for resolving the e_2
func MyServeFile(w http.ResponseWriter, r *http.Request, fileName string) {
	//var fileName = "html/aboutAuthor_ru.html"
	var b []byte
	var msg string
	var err error
	if b, err = ioutil.ReadFile(fileName); err != nil {
		panic(fmt.Sprintf("MyServeFile: ioutil.ReadFile err=%s", err.Error()))
	}
	msg = fmt.Sprintf("<p>%s</p>", b)
	w.WriteHeader(200)
	w.Write([]byte(msg))
	//if isDebug(serverMode) {
	//	fmt.Printf("MyServeFile: RequestURI =%v, actualLang=%v, file=%s\n%s\n", r.RequestURI, getLang(r), fileName, b[:100])
	//}
}

//220128 10:29 220421 21:19
func isDebug(mode int) bool {
	return digInPos(mode, 1) > 0
}
func isHTTPS(mode int) bool {
	return digInPos(mode, 2) > 0
}
func isSms(mode int) bool {
	return digInPos(mode, 3) > 0
}
func isHidingOut(mode int) bool {
	//fmt.Printf("---------digInPos:---mode=%d, pos=%d, res=%d\n", mode, pos, res)
	return digInPos(mode, 4) > 0
}

//------------------------
//220421 21:18 see Progects/golang/220421_lastdig
func digInPos(mode int, pos int) int {
	if !(pos > 0) {
		panic("digInPos: pos<1")
	}
	if mode < 0 {
		panic("digInPos: mod<0")
	}
	var power10 = func(n int, p int) int { // return n^10
		n = 1
		for i := 0; i < p; i++ {
			n = n * 10
		}
		return n
	}
	var res = (mode / power10(1, pos-1)) % 10 // in begin it gets the remainder, next it gets the last gigit of the remaider
	//fmt.Printf("---------digInPos:---mode=%d, pos=%d, res=%d\n", mode, pos, res)
	return res
}

//220422 15:24 for printProblem
func debugLevel() int {
	return digInPos(serverMode, 1)
}

func normSms(sms string) string {
	const maxSmsLen = 20
	var bound int
	if len(sms) <= maxSmsLen {
		bound = len(sms)
	} else {
		bound = maxSmsLen
	}
	return strings.ReplaceAll(sms[:bound], " ", "_")
}

//220127 13:36
//_______15:58
func sendSms(sms, to string) (answer string, err error) {
	//from https://sms.ru/api/send
	//https://sms.ru/sms/send?api_id=C30C93EC-C295-9E8F-6B96-6C589E55D8D1&to=79536130260,74993221627&msg=hello+world&json=1
	const api_id = "C30C93EC-C295-9E8F-6B96-6C589E55D8D1"
	var resp *http.Response
	var b []byte = make([]byte, 2000)

	if to == "" {
		to = "9536130260"
	}
	sms = normSms(sms)
	var URI = fmt.Sprintf("https://sms.ru/sms/send?api_id=%s&to=%s&msg=%s&json=1", api_id, to, sms)
	if resp, err = http.Get(URI); err != nil {
		answer = ""
		err = fmt.Errorf("sendSms: http.Get(URI) err=%s", err.Error())
		return
	}
	resp.Body.Read(b)
	resp.Body.Close()
	answer = string(b)
	if resp.StatusCode != 200 {
		err = fmt.Errorf("sendSms: http.Get(URI) %s", answer)
		answer = ""
		return
	}
	return
}

//220128 06:34
//It leaves no more than max runes in s
//if max<=0 it panics
func truncStr(s string, max int) string {
	var rCount int
	var res string
	if max <= 0 {
		panic(fmt.Sprintf("truncStr: not permitted value of max=%d", max))
	}
	if strings.Count(s, "") < max-1 {
		return s
	}
	for _, r := range s {
		if rCount > max {
			break
		}
		res = res + string(r)
		rCount++
	}
	return res
}

//220128 08:41 Returns list of files sorted by modification time (the later in begining)
func getLetList() []os.FileInfo {
	var err error
	var letDir *os.File
	var letInfos []os.FileInfo
	if letDir, err = os.Open("letters"); err != nil {
		panic(fmt.Sprintf("getLetList: os.Open(\"letters\" err=%s)", err.Error()))
	}

	if letInfos, err = letDir.Readdir(0); err != nil {
		panic(fmt.Sprintf("getLetList: letDir.Readdir(0) err=%s)", err.Error()))
	}
	letDir.Close()
	/*
		fmt.Println("truncLetDir================")
		for i := 0; i < len(letInfos); i++ {
			fmt.Printf("%s----%s\n", letInfos[i].Name(), letInfos[i].ModTime().Format(timeFormat))
		}
		fmt.Println("================truncLetDir")
	*/
	sort.SliceStable(letInfos, func(i, j int) bool {
		return letInfos[i].ModTime().Format(timeFormat) < letInfos[j].ModTime().Format(timeFormat)
	})
	/*
		fmt.Println("truncLetDir-----------------")
		for i := 0; i < len(letInfos); i++ {
			fmt.Printf("%s----%s\n", letInfos[i].Name(), letInfos[i].ModTime().Format(timeFormat))
		}
		fmt.Println("----------------truncLetDir")
	*/
	return letInfos
}
func truncLetDir() {
	var letInfos []os.FileInfo = getLetList()

	if len(letInfos) <= maxLetters {
		return
	}
	letInfos = letInfos[:maxLetters-1]

	for _, item := range letInfos {
		os.Remove("letters/" + item.Name())
	}
}

//220302 17:30 for saveAccountList (accounts.go)
//220307 07:37 An example of result: "lang=en qqq=", that is the qqq key is empty.
//220404 09:38 see var accountCompsSpr(="-;;-") var optionsSpr(=";--;") var optionKVspr(=";==;")
//220419 13:26 If whatKey=="" it gives all as "key;==;val;--;key;==;val ..."
//If whatKey!="" it gives "key;==;val" where key==whatKey
func optsToStr(opts map[string]string, whatKey string) string {
	var res string
	var err error
	var valAsStr string
	var keyCount int //220307 09:23

	if opts == nil { //220307 08:58
		panic("optsToStr: nill not allowed as parameter")
	}
	if len(opts) < minOptions {
		panic(fmt.Sprintf("utils.go>optsToStr:len(opts) < %d", minOptions))
	}

	for key, val := range opts {
		if err = goodString(key); err != nil {
			panic("optToStr: Bad key:" + err.Error())
		}
		if key == "" {
			panic(fmt.Sprintf("utils.go>optToStr: empty key; keyCount=%d", keyCount))
		}
		if err = goodString(val); err != nil {
			panic("optToStr: bad value:" + err.Error())
		}
		if whatKey != "" {
			if key != whatKey {
				continue
			}
		}

		valAsStr = fmt.Sprintf("%s%s%s", key, optionKVspr, val)

		keyCount++
		if keyCount == 1 {
			res = valAsStr
		} else {
			res = res + optionsSpr + valAsStr
		}
	}
	return res
}

//220303 08:14 --16:18
func goodString(s string) (err error) {
	if s == "" {
		return
	}
	for _, val := range s { //220307 10:48 11:14 - it is seemed useless but let it be
		if val == '=' {
			err = fmt.Errorf("%s is bad (= character is not allowed)", s)
			return
		}

	}
	if !utf8.Valid([]byte(s)) {
		err = fmt.Errorf("%s is bad (as utf8.Valid says)", s)
		return
	}
	for ind, runeVal := range s {
		if runeVal <= 20 {
			return fmt.Errorf("%s has bad run (%s) at %i", s, string(runeVal), ind)
		}
	}
	return
}

//220307 07:33 for saveAccountList (accounts.go)
//220310 11:27
//220404 09:46 see var accountCompsSpr(="-;;-") var optionsSpr=(";--;")
func strToOpts(s string) (opts map[string]string) {
	var optsSlice, optSlice []string

	if s == "" {
		panic("strToOpts: An empty parameter is not allowed")
	}

	opts = make(map[string]string)
	optsSlice = strings.Split(s, optionsSpr)
	//fmt.Printf("-----utils.go>strToOpts: optsSlice=%v\n", optsSlice)
	if len(optsSlice) < minOptions {
		panic(fmt.Sprintf("Bad record of option(%s): not enough options, minOptions=%d", minOptions))
	}
	//fmt.Printf("----utils.go>strToOpts: s=%s\n", s)
	for ind, keyVal := range optsSlice {
		optSlice = strings.Split(keyVal, optionKVspr)
		if len(optSlice) != 2 { //220307 11:08 - it is a need to check that splitting "qqq=" gives a slice with two components
			//220404 10:52 the optSlice must contain exactly two components
			panic(fmt.Sprintf("Bad record of option(%s)(index%d): no two componenct (key and value)(optionKVspr=%s)", keyVal, ind, optionKVspr))
		}
		//fmt.Printf("----utils.go>strToOpts: optSlice=%v\n", optSlice)
		opts[optSlice[0]] = optSlice[1]
	}
	return opts
}

//220309 10:52 it panics if n is not a valid account name
func checkAccontName(n string) string {
	if n == "" {
		panic("checkAccontName: name cannot be empty.")
	}
	return n
}

//220309 10:52 it panics if tp is not represent a valid account.Tp
//220404 14:16
func convAccontTp(tp string) (res int) {
	var err error
	if res, err = strconv.Atoi(tp); err != nil {
		panic(fmt.Sprintf("convAccontTp: %s is not a valid integer", tp))
	}

	//if res != 0 {
	//	panic(fmt.Sprintf("convAccontTp: only 0 is allowed but there is %d", res))
	//}

	return res
}

//220309 14:00 it panics if tm is not represent a valid account.RegTm
func convAccountTm(tm string) (res time.Time) {
	var err error
	if res, err = time.Parse(timeFormat, tm); err != nil {
		panic(fmt.Sprintf("convAccontTm: conversion %s err=%s", tm, err.Error()))
	}
	return res
}

// 220421 21:57
func printDebug(msg string) {
	if !isDebug(serverMode) {
		return
	}
	fmt.Printf("DEBUG-----%s\n", msg)
}

func printProblem(problem, msg string) {
	var caller string
	var file string
	var line int
	if !isDebug(serverMode) {
		return
	}
	_, file, line, _ = runtime.Caller(1)
	caller = fmt.Sprintf("%s(%d)", file, line)
	fmt.Printf("%s-----%s;caller=%s\n", problem, msg, caller)
	//time.Sleep(time.Second)
}

//22030404 07:25
//220405 10:11
func removeOldLogs() {
	//func removeOldLogs() error {
	//var cmd = exec.Command("rm", "*.log")
	//return cmd.Run()
	files, err := filepath.Glob("*.log")
	if err != nil {
		panic(fmt.Sprintf("utils.go>removeOldLogs: getting file list err=%s", err.Error()))
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(fmt.Sprintf("utils.go>removeOldLogs: removing file=%s  err=%s", f, err.Error()))
		} else {
			fmt.Printf("utils.go>removeOldLogs: %s was removed\n", f)
		}
	}
}

//220404 11:53
func byteSliceToStrRepresentation(bs []byte) (res string) {
	res = fmt.Sprintf("%v", bs)
	return
}

//220404 13:00
//It is the inverse of byteSliceToStrRepresentation
//That is it expects s as "[1 2 3 4 ...]"
func byteStrRepresentationToByteSlice(s string) (res []byte) {
	for _, member := range []byte(s) {
		res = append(res, member)
	}

	return
}

//220407 14:27 I did not find any better than to write this utilite.
func writeAllToFile(f *os.File, buf []byte) {
	var bytes int = len(buf)
	var totalWriten int
	var writen int
	var maxCount = 10
	var count int
	var err error

	if bytes == 0 {
		panic("writeAllToWile: nothing to write")
	}

	for totalWriten != bytes {
		count++
		writen, err = f.Write(buf)
		if err != nil {
			panic(fmt.Sprintf("writeAllToWile: writing err=%s", err.Error()))
		}
		if count > maxCount {
			panic(fmt.Sprintf("writeAllToWile: count more that %d", maxCount))
		}
		totalWriten = totalWriten + writen
	}
}

//220420 11:53 For prolongeAccount
func compareSlices(sl1, sl2 []byte) (bool, error) {
	var l1 = len(sl1)
	var l2 = len(sl2)
	if l1 == 0 || l2 == 0 {
		panic("compareSlices: byte slices must not be of zero length")
	}
	if l1 != l2 {
		return false, fmt.Errorf("lengths is not equal")
	}

	for i := 0; i < l1; i++ {
		if sl1[i] != sl2[i] {
			return false, fmt.Errorf("In position %d is not equality", i)
		}
	}
	return true, nil
}

func haltAll(msg string) {
	fmt.Printf("The program has been halted at %v with:%s", time.Now(), msg)
	os.Exit(111)
}
