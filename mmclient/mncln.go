//220418 07:42
//mich412@mich412-A320M-S2H-V2:~/Progects/freelancer/mmst/mmclient$ ./mmclient
//-----get: uri=https://www.dost346.ru:8080/qGet: err=Get "https://www.dost346.ru:8080/q": x509: certificate signed by unknown authority;uri=https://www.dost346.ru:8080/q
//So, from overview of  https://pkg.go.dev/net/http :
//For control over proxies, TLS configuration, keep-alives, compression, and other settings, create a Transport: ...
//Also from https://stackoverflow.com/questions/38822764/how-to-send-a-https-request-with-a-certificate-golang; second answer
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
	"log"
)

var defurl string = "https://www.dost346.ru:8080/q" // default value if there is not an argument
var defmeth = "get"

var help = `
220416 09:17
It is http(s) CLI client that regards the answer as a plain text and print it to stdout
A command line is :
prog urg1, urg2, ... urgN
where
prog is a reference to this programm, for example ./mmclient
urg<number> is a command line agrument with format of <option name>=<option value>, for example URL="https://dost346.ru:8080/qqq?xxx=y zzz"
As usual, the option value must be enslosed by quotes if it contains blank characters.
Arguments are not mandatory. That is for each there is a default value.

There are options (the case is significant):
h - it is exclusive option. If it is in command line the prog ignoring all and only print this help text.
URI - uri, by which the prog answers a http(s) server, for example URL="https://www.dost346.ru:8080/qqq?xxx=y zzz";
default value - https://dost346.ru:8080/ind
METH - a method of a request, for example METH=post;
default value - get

If in the working directory of the prog there is file "options.txt" the prog prefer optionss from it than from the command line.
That is the prog reading its option from the file (excluding the option "h").

RETURNS:
if there is option "h" (with any value) the prog returns 0 and prints this text.
Otherwise it questions a server;
if is any answer from the server it returns 0 and print the answer, overwise it returns 1 and primts error message.
`

func printResp(resp *http.Response) {
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
}

//220416 09:42 It is from the mmst. This scheme seems very good.
//That is it allows arbitrary amount of arguments and only one format of a command line.
func setArgs() {
	for i := 1; i < len(os.Args); i++ {
		var splitedArg = strings.Split(os.Args[i], "=")
		if len(splitedArg) != 2 {
			fmt.Printf("It are allowed arguments with format name=value, but is %v\n", os.Args[i])
			os.Exit(1)
		}
		switch splitedArg[0] {
		case "h":
			fmt.Println(help)
			os.Exit(0)
		case "URI":
			defurl = splitedArg[1]
		case "METH":
			if splitedArg[1] != defmeth {
				fmt.Printf("It is allowed get method, but is %v\n", splitedArg[0])
				os.Exit(1)
			} else {
				defmeth = splitedArg[1]
			}
		default:
			fmt.Printf("There is not allowed option %v\n", splitedArg[0])
			os.Exit(1)
		}
	}
}

func main() {
	//var err error
	//var res string
	setArgs() //if there is "h" option that the prog will halt

	//220418 08:12
	caCert, err := ioutil.ReadFile("cert")
	if err != nil {
		log.Fatal(err)
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
	if resp, err := client.Get(defurl); err != nil {
		fmt.Printf("Getting %s err=%s\n", defurl, err.Error())
		os.Exit(1)
	} else {
		printResp(resp)
	}
	//-------------------------------------------------

}
