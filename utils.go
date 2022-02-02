// utils
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
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
	//var err error
	const MaxMode = 111
	for i := 1; i < len(os.Args); i++ {
		var splitedArg = strings.Split(os.Args[i], "=")
		if len(splitedArg) != 2 {
			fmt.Printf("It are allowed arguments with format name=value, but is %v\n", os.Args[i])
			os.Exit(1)
		}
		switch splitedArg[0] {
		//		case "pg":
		//			PG_CONN_STR = splitedArg[1]
		case "la":
			listeningAddr = splitedArg[1]
			break
		case "pw":
			passWord = splitedArg[1]
			break
		case "mode":
			serverMode, _ = strconv.Atoi(splitedArg[1])
			if serverMode < 0 || serverMode > MaxMode {
				serverMode = 0
			}
			break
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

//220128 10:29
func isDebug(mode int) bool {
	return mode%10 != 0
}
func isHTTPS(mode int) bool {
	var withoutLastdigit int
	withoutLastdigit = mode / 10
	return withoutLastdigit%10 != 0
}
func isSms(mode int) bool {
	var withoutTwoLastgit int
	withoutTwoLastgit = mode / 100
	return withoutTwoLastgit%10 != 0
}

//------------------------

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
