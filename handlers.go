// handlers
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/user"
	"strconv"

	"html/template"
	"io/ioutil"
	"strings"
	"time"

	"bufio"
)

var buffs [][]byte

//220113 13:02
func historyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(history))
}

//220118 07:31 220331 11:54
//220405 10:45
func accountsHandler(w http.ResponseWriter, r *http.Request) {
	var whatField, whatOpt string
	if r.FormValue("pw") != passWord {
		msg := fmt.Sprintf("%s > bad password", getRequestBrief(r))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(500)
		w.Write([]byte(msg))
		return
	}

	whatField = r.FormValue("field")
	whatOpt = r.FormValue("opt")

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(getAccountsAsHTML(whatField, whatOpt)))
}

//220113 04:43 Only sense of this to resolve the enigma of e_2 (see the history of 220113 04:17)
func e_2Handler(w http.ResponseWriter, r *http.Request) {
	var fileName = "html/aboutAuthor_ru.html"
	var b []byte
	var msg string
	var err error
	if b, err = ioutil.ReadFile(fileName); err != nil {
		panic(fmt.Sprintf("e_2Handler: ioutil.ReadFile err=%s", err.Error()))
	}
	msg = fmt.Sprintf("<p>%s</p>", b)
	w.WriteHeader(200)
	w.Write([]byte(msg))
	//fmt.Printf("e_2Handler RequestURI =%v, actualLang=%v\n", r.RequestURI, lang)

}

//220112 07:54
func myFamilyHandler(w http.ResponseWriter, r *http.Request) {
	//var fileName string
	//fmt.Printf("aboutAuthorHandler RequestURI =%v, actualLang=%v, file=%v\n", r.RequestURI, lang, fileName)
	http.ServeFile(w, r, "image/myFamily.png")

}

//220202 15:37
func cssHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "mystyle.css")
}

//220111 14:59; 220125 05:50
func aboutAuthorHandler(w http.ResponseWriter, r *http.Request) {
	var fileName__ string
	if getLang(r) == "en" {
		fileName__ = "html/aboutAuthor_en.html"
	} else {
		fileName__ = "html/aboutAuthor_ru.html"
	}
	//http.ServeFile(w, r, fileName__)
	MyServeFile(w, r, fileName__)
	if isDebug(serverMode) {
		fmt.Printf("aboutAuthorHandler: lang=%s; file=%s\n", getLang(r), fileName__)
	}
}

//220111 05:27
func mmsite_scriptHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "mmsit.js")
}

//220118 16:51
//220124 17:27 See also langErr220124; On 200 it returns "ru" or "en"
func changeLangHandler(w http.ResponseWriter, r *http.Request) {
	var newLang string = r.URL.RawQuery
	var msg string
	var opts map[string]string
	//var accName string
	if newLang == "" {
		msg = fmt.Sprintf("changeLangHandler: RequestURI =%v, error: no language to change", r.RequestURI)
		w.WriteHeader(400)
		w.Write([]byte(msg))
		return
	}
	if newLang != "en" && newLang != "ru" {
		msg = fmt.Sprintf("changeLangHandler: RequestURI =%v, bad language to change : %s", r.RequestURI, newLang)
		w.WriteHeader(400)
		w.Write([]byte(msg))
		return
	}
	opts = getOptions(r)
	opts["lang"] = newLang
	setOptions([]byte(accountName(r)), opts)
	msg = fmt.Sprintf("%s", newLang)
	w.WriteHeader(200)
	w.Write([]byte(msg))
}

//210111 05:50
//220419 18:05
func indHandler(w http.ResponseWriter, r *http.Request) {
	type templData struct {
		ReqNum int64
		Lang   string
	}
	var err error
	var fileName string
	var templ *template.Template

	var tD templData
	tD.ReqNum = flr.feelerCount
	tD.Lang = getLang(r)

	printProblem("problem_220415", fmt.Sprintf("indHandler: tD.Lang=%s", tD.Lang))

	//220131 04:36 This nonsense has existed long ago and it had not been noticed up to the last friday.
	// This phenomenon is such interesting that it deserve a name: langErr220131

	if tD.Lang == "en" { // 220419 18:24 And here I have span as a louse on a comb all day
		fileName = "./html/ind_en.html"
	} else {
		fileName = "./html/ind_ru.html"
	}
	/* 220415 09:43 This stratum is unnecessary as all panics are intercepted by the feeler.
	defer func() {
		if rec := recover(); rec != nil {
			//panicMessage := fmt.Sprintf("(Addr=%v;N=%v) panic:%v", r.RemoteAddr, getRequestCounter(), rec)
			panicMessage := fmt.Sprintf("%v", rec)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			w.Write([]byte(panicMessage))
			fmt.Println("indHandler panics:", panicMessage)
		}
	}()
	*/
	if templ, err = template.ParseFiles(fileName); err != nil {
		panic(fmt.Sprintf("indHandler: parsing %s err=%s", fileName, err.Error()))
	}

	if err = templ.Execute(w, tD); err != nil {
		panic(fmt.Sprintf("indHandler: executing %s err=%s", fileName, err.Error()))
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

}

func indHandler_220108(w http.ResponseWriter, r *http.Request) {
	type templData struct {
		ReqNum int64
	}
	var err error
	var fileName = "ind.html"
	var templ *template.Template

	var tD templData
	tD.ReqNum = flr.feelerCount

	defer func() {
		if rec := recover(); rec != nil {
			panicMessage := fmt.Sprintf("(Addr=%v;N=%v) panic:%v", r.RemoteAddr, getRequestCounter(), rec)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			w.Write([]byte(panicMessage))
			fmt.Println("!!!!!!!indHandler panics:", panicMessage)
		}
	}()

	if templ, err = template.ParseFiles(fileName); err != nil {
		panic(fmt.Sprintf("indHandler_220108: parsing %s err=%s", fileName, err.Error()))
	}

	if err = templ.Execute(w, tD); err != nil {
		panic(fmt.Sprintf("indHandler_220108: executing %s err=%s", fileName, err.Error()))
	}
}

func showMemStatHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var cmLen int //сm stands for "consume memory"
	var cmLenStr string
	var cm []byte
	var memRep string
	var start = time.Now()
	defer func() {
		if rec := recover(); rec != nil {
			panicMessage := fmt.Sprintf("(Addr=%v) panic:%v", r.RemoteAddr, rec)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			w.Write([]byte(panicMessage))
			fmt.Println("!!!!!!!showMemStat panics:", panicMessage)
		}
	}()

	cmLenStr = r.FormValue("cm")
	if cmLenStr != "" {
		if cmLen, err = strconv.Atoi(cmLenStr); err != nil {
			panic("There is cm but it cannot be converted to int")
		}
		if cmLen < 1000 {
			panic("There is cm<1000")
		}
		cm = make([]byte, cmLen)
		buffs = append(buffs, cm)
	}

	memRep = memStatStr(start)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(memRep))

	fmt.Printf("(%v)showMemStatHandler\n", getRequestCounter())
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			panicMessage := fmt.Sprintf("(Addr=%v) panic:%v", r.RemoteAddr, rec)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			w.Write([]byte(panicMessage))
			fmt.Println("!!!!!!!stopHandler panics:", panicMessage)
		}
	}()

	fmt.Printf("(%v)stopHandler\n", getRequestCounter())
	os.Exit(0)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	var usr *user.User
	var hostName string
	var err error
	var msg string
	defer func() {
		if rec := recover(); rec != nil {
			panicMessage := fmt.Sprintf("(Addr=%v) panic:%v", r.RemoteAddr, rec)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			w.Write([]byte(panicMessage))
			fmt.Println("!!!!!!!aboutHandler panics:", panicMessage)
		}
	}()

	if r.FormValue("pw") != passWord {
		panic(fmt.Sprintf("aboutHandler panic: wrong password"))
	}

	if usr, err = user.Current(); err != nil {
		panic(fmt.Sprintf("aboutHandler panic of getting user info:%v", err.Error()))
	}
	if hostName, err = os.Hostname(); err != nil {
		panic(fmt.Sprintf("aboutHandler panic of getting host name:%v", err.Error()))
	}
	msg = fmt.Sprintf("<p>Uid=%v,<br> login==%v,<br> name=%v,<br> HomeDir=%v,<br>Host=%v</p>",
		usr.Uid, usr.Username, usr.Name, usr.HomeDir, hostName)

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(msg))
	//fmt.Printf("(%v)aboutHandler\n", getRequestCounter())

}

//220108 07:34
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "favicon.ico")
}

//220127 06:22
//sms?text=...
//220127 14:18
func smsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var Sms string
	var answer string

	if r.FormValue("pw") != passWord {
		panic(fmt.Sprintf("smsHandler panic: wrong password"))
	}

	Sms = r.FormValue("sms")

	if answer, err = sendSms(Sms, ""); err != nil {
		answer = fmt.Sprintf("smsHandler: sendSms err=%s", err.Error())
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(answer))
}

//220127 15:24 220128 04:50
func letterHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var cont string //the content of the letter
	var fileName = fmt.Sprintf("letters/letter_%s_%s.txt", accountName2(r), time.Now().Format("20060102_15:04:05"))
	var letFile *os.File
	var letVolume int
	var msgSms string
	var msgLet string

	if cont = r.FormValue("cont"); cont == "" {
		panic(fmt.Sprintf("letterHamdler: empty letter"))
	}

	cont = truncStr(cont, maxChars)

	letVolume = strings.Count(cont, "")

	if letFile, err = os.Create(fileName); err != nil {
		panic(fmt.Sprintf("letterHamdler: os.Create err=%s", err.Error()))
	}

	if _, err = letFile.WriteString(cont); err != nil {
		panic(fmt.Sprintf("letterHamdler: letFile.WriteString err=%s", err.Error()))
	}

	if isSms(serverMode) {
		msgSms = fmt.Sprintf("let from %s", accountName(r))
		if _, err = sendSms(msgSms, ""); err != nil {
			panic(fmt.Sprintf("letterHamdler: sendSms err=%s", err.Error()))
		}
		msgLet = fmt.Sprintf("The letter of %d (from %s) was saved successfully with a name of %s WITH sending sms", letVolume, accountName2(r), fileName)
	} else {
		msgLet = fmt.Sprintf("The letter of %d (from %s) was saved successfully with a name of %s WITHOUT sms", letVolume, accountName2(r), fileName)

	}

	truncLetDir()

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(msgLet))
}

//220330 09:16
//func registraterHandler(w http.ResponseWriter, r *http.Request) {
//	var err error
//	w.Header().Add("Content-Type", "text/html; charset=utf-8")
//	w.WriteHeader(200)
//	w.Write([]byte(msgLet))
//}

//220331 08:12
func myAccountHandler(w http.ResponseWriter, r *http.Request) {
	//var err error
	var msg, opts string
	var acc *Account
	var accName = r.Context().Value(AccNameCtxKey).([]byte)
	//if accName == "?" {//220422 06:29
	//	panic("myAccountHandler: no account name")
	//}
	if acc = getAccount([]byte(accName)); acc == nil {
		panic(fmt.Sprintf("myAccountHandler: no account of name %s", accName))
	}
	msg = fmt.Sprintf("<h2>Accunt %v or(in hex coding)\"%x\" </h2>", accName, string(accName))
	for key, val := range acc.Options {
		opts = opts + fmt.Sprintf("<p>%s=%s</p>", key, val)
	}
	msg = msg + opts + "<p>-------------------------------</p>"
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(msg))
}

//220407 12:13
//A crossroads: 1. to use f.log (*os.File) 2. to use ioutil.ReadFile(fileName)
//12:23 I select second way; It seems to me more educational
func showFeelerLogHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var line string
	var substr string
	var f *os.File
	var sc *bufio.Scanner

	//printDebug("showFeelerLogHandler: srart")

	if f, err = os.Open(flr.logFileName); err != nil {
		panic(fmt.Sprintf("showFeelerLogHandler: open front log err=%s", err.Error()))
	}
	substr = r.FormValue("substr")
	sc = bufio.NewScanner(f)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	for sc.Scan() {
		line = sc.Text()
		//printDebug(fmt.Sprintf(""))
		if strings.Contains(line, substr) {
			line = line + "\n"
			w.Write([]byte(line))
		}
	}
	if err = sc.Err(); err != nil {
		fmt.Fprintln(w, "Scanning front log err=", err)
	}
	fmt.Fprintln(w, "--------------------------")

}

//2200408 10:06
func helpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(help))
}

//220412 15:56
func showGeneralLogHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var line string
	var substr string
	var f *os.File
	var sc *bufio.Scanner

	//printDebug("showFeelerLogHandler: srart")

	if f, err = os.Open(generalLogFileName); err != nil {
		panic(fmt.Sprintf("showGeneralLogHandler: open log err=%s", err.Error()))
	}
	substr = r.FormValue("substr")
	sc = bufio.NewScanner(f)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)

	for sc.Scan() {
		line = sc.Text()
		//printDebug(fmt.Sprintf(""))
		if strings.Contains(line, substr) {
			line = line + "\n"
			w.Write([]byte(line))
		}
	}
	if err = sc.Err(); err != nil {
		fmt.Fprintln(w, "Scanning general  log err=", err)
	}
	fmt.Fprintln(w, "--------------------------")

}

//220422 14:48 Especially for resolving the problem_220415
func problem_220415Handler(w http.ResponseWriter, r *http.Request) {
	type templData struct {
		ReqNum int64
		Lang   string
	}
	var err error
	var fileName = "problem_220415.html"
	var templ *template.Template

	var tD templData
	tD.ReqNum = flr.feelerCount
	tD.Lang = getLang(r)

	if templ, err = template.ParseFiles(fileName); err != nil {
		panic(fmt.Sprintf("indHandler_220108: parsing %s err=%s", fileName, err.Error()))
	}

	if err = templ.Execute(w, tD); err != nil {
		panic(fmt.Sprintf("indHandler_220108: executing %s err=%s", fileName, err.Error()))
	}
}

//2200425 17:57
func longOperHandler(w http.ResponseWriter, r *http.Request) {
	var msg string
	var start = time.Now()
	time.Sleep(3 * time.Second)
	var dur = time.Since(start)
	msg = fmt.Sprintf("The long operation here, dur=%v", dur)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(msg))
}
