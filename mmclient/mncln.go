//220418 07:42
//mich412@mich412-A320M-S2H-V2:~/Progects/freelancer/mmst/mmclient$ ./mmclient
//-----get: uri=https://www.dost346.ru:8080/qGet: err=Get "https://www.dost346.ru:8080/q": x509: certificate signed by unknown authority;uri=https://www.dost346.ru:8080/q
//So, from overview of  https://pkg.go.dev/net/http :
//For control over proxies, TLS configuration, keep-alives, compression, and other settings, create a Transport: ...
//Also from https://stackoverflow.com/questions/38822764/how-to-send-a-https-request-with-a-certificate-golang; second answer
//220418 13:46 if the func setArgs() returns false the func setAgrsFromFile() will be invoked
//220420 09:26 Adding the time measuring and versioning
package main

import (
	"fmt"
	"net/http"

	"bytes"
	"strings"

	"os"

	"time"

	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"bufio"
)

const OptionsFile = "options.txt"

var defurl string = "https://www.dost346.ru:8080/q" // default value if there is not an argument
var defmeth = "get"
var defcookies string //a default is absent of cookies

var help = `
220416 09:17 220419 07:33
It is http(s) CLI client that makes a http(s) requst and prints result to to stdout
A command lines are :
prog
prog h
where "prog" is a call of this programm, for example ./mmclient
If there is the argument "h" (that is when "prog h") the program print (to console) this text and ended its work with code 0.
Otherwise the program demands existing in its working directory a file of "options.txt" from which it takes its options.

There are options (the case is significant):
URI - uri, by which the prog answers a http(s) server, for example URL="https://www.dost346.ru:8080/qqq?xxx=y zzz";
METH - a method of a request, for example METH=get;
COOKIES - a cookies that will sent with a requst
RETURNS:
If there is option "h" the prog returns 0
Otherwise it returns 0 if sending a request takes place and 1 if no request was sent
In any case the program print a mesage to console.
`

//220420 09:23 Adding the time measuring
func printResp(resp *http.Response, start time.Time) {
	var err error
	var buff *bytes.Buffer

	fmt.Printf("--------------Response for %s\n", defurl)
	buff = bytes.NewBuffer(nil)
	if err = resp.Write(buff); err != nil {
		fmt.Printf("writing to buff err=%s\n", err.Error())
	} else {
		fmt.Printf("%s\n", buff.String())
	}
	resp.Body.Close()

	fmt.Printf("----------------dur=%s\n", time.Since(start).String())
}

//220419 08:52
func toTwo(sl []string, link string) []string {
	var res []string = make([]string, 2)
	var rest string
	switch len(sl) {
	case 0, 1:
		panic("toTwo: sourse slice len = 0 or 1")
	case 2:
		return sl
	}
	for i := 1; i < len(sl); i++ {
		if i == 1 {
			rest = rest + sl[i]
		} else {
			rest = rest + link + sl[i]
		}
	}
	res[0] = sl[0]
	res[1] = rest
	return res
}

//220418 13:46
func setAgrsFromFile() {
	var err error
	var line string
	var f *os.File
	var sc *bufio.Scanner

	var optSlice []string

	//printDebug("showFeelerLogHandler: srart")

	if f, err = os.Open(OptionsFile); err != nil {
		panic(fmt.Sprintf("setAgrsFromFile: open %s err=%s", OptionsFile, err.Error()))
	}
	sc = bufio.NewScanner(f)
	for sc.Scan() {
		line = sc.Text()
		if (len(line) != 0) && (line[0] != '#') {
			optSlice = strings.Split(line, "=")
			if len(optSlice) == 1 {
				panic(fmt.Sprintf("setAgrsFromFile: line %s is not comment and has not = character", line))
			}
			if len(optSlice) > 2 {
				//panic(fmt.Sprintf("setAgrsFromFile: line %s has more than one = character", line))
				optSlice = toTwo(optSlice, "=")
			}
			switch optSlice[0] {
			case "URI":
				defurl = optSlice[1]
			case "COOKIES":
				defcookies = optSlice[1]
			default:
				panic(fmt.Sprintf("setAgrsFromFile: line %s has illegal option %s", line, optSlice[0]))
			}
			//fmt.Println("--------" + optSlice[1])

		}
	}
	if err = sc.Err(); err != nil {
		panic(fmt.Sprintf("setAgrsFromFile: scanning %s", err.Error()))
	}

}

func main() {
	fmt.Println("mmclient 1.0; For help invoke with argument h")
	if len(os.Args) > 1 { //if there is "h" option that the prog will halt
		if os.Args[1] != "h" {
			panic("Allowed only h argument")
		} else {
			fmt.Println(help)
			os.Exit(0)
		}
	}
	setAgrsFromFile()

	//220418 08:12
	caCert, err := ioutil.ReadFile("cert")
	if err != nil {
		fmt.Printf("Readig cart file err=%s\n", err.Error())
		os.Exit(1)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			RootCAs: caCertPool,
		},
	}
	client := &http.Client{Transport: tr}
	//220419 06:19
	req, err := http.NewRequest("GET", defurl, nil)
	if err != nil {
		fmt.Printf("http.NewRequest (for %s) err=%s\n", defurl, err.Error())
		os.Exit(1)
	}
	if defcookies != "" {
		req.Header.Add("cookie", defcookies)
	}

	var start = time.Now()
	if resp, err := client.Do(req); err != nil {
		fmt.Printf("client.Do (for %s) err=%s\n", defurl, err.Error())
		os.Exit(1)
	} else {
		printResp(resp, start)
	}
	//-------------------------------------------------

}
