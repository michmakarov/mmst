//220503 07:43 outsess
//Obviously that the range of requests must be taken without the registration
package main

import (
	"fmt"
	"net/http"
)

func qqq() {
	fmt.Println("Hello World!")
}

//220503 13:09 True is means that the r queries a page that must be answered without a client registration and it has been answered
//_______19:36
func reqAnswered(w http.ResponseWriter, r *http.Request, accName []byte, accRes byte) (res bool) {
	var path = r.URL.Path
	res = true // if there are not a target path the default case will establish it in false.
	switch path {
	case "/q":
		s := fmt.Sprintf("There is /q debug request; accName=%v; accres=%d; RA=%s", accName, accRes, r.RemoteAddr)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(s))
	case "/":
		indHandler(w, r)
	case "/help":
		helpHandler(w, r)
	case "/history":
		historyHandler(w, r)
	case "/longOper":
		longOperHandler(w, r)
	case "/mmsite_script":
		mmsite_scriptHandler(w, r)
	case "/css":
		cssHandler(w, r)
	case "/favicon.ico":
		faviconHandler(w, r)
	default:
		res = false
	}

	return
}
