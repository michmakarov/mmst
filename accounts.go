// 220117 18:26 accounts
package main

import (
	"container/list"
	"fmt"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

type Account struct {
	Name    string            // identificator
	Tp      int               //type 0 - automatically created; the name is remote address; no password
	Options map[string]string //options, sequence of substring of format key=value,for example lang=en somekey=somevalue ...
	RegTm   time.Time         //Time of regidtrarion
}

var accounts *list.List
var accountsMtx sync.Mutex // it queues queries for access to the accounts

func init() {
	if accounts == nil {
		accounts = list.New()
	}
}

//220118 05:43 The func takes an account name and returns a pointer to an account if it exists or nil if not
func getAccount(accName string) *Account {
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	for e := accounts.Front(); e != nil; e = e.Next() {
		if e.Value.(*Account).Name == accName {
			return e.Value.(*Account)
		}
	}
	return nil
}

//220118 17:12
func accountName(r *http.Request) string {
	var accName string
	var err error
	if accName, _, err = net.SplitHostPort(r.RemoteAddr); err != nil {
		panic(fmt.Sprintf("accccountName: net.SplitHostPort err=%s", err.Error()))
	}
	return accName
}

//220118 16:22 The func panics if opts==nil
func setOptions(accName string, opts map[string]string) {
	if opts == nil {
		panic(fmt.Sprintf(": no account no opt (nil) for %s", accName))
	}
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	for e := accounts.Front(); e != nil; e = e.Next() {
		if e.Value.(*Account).Name == accName {
			e.Value.(*Account).Options = opts
		}
	}

}

//220118 15:04 The func returns options related a given request.
//If the request is not correspond some account the func panics.
func getOptions(r *http.Request) map[string]string {
	var err error
	var accName string
	if accName, _, err = net.SplitHostPort(r.RemoteAddr); err != nil {
		panic(fmt.Sprintf("getCurrOptions: net.SplitHostPort err=%s", err.Error()))
	}
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	for e := accounts.Front(); e != nil; e = e.Next() {
		if e.Value.(*Account).Name == accName {
			return e.Value.(*Account).Options
		}
	}
	panic(fmt.Sprintf("getCurrOptions: no account for %s", r.RequestURI))
}

//220118 15:39 The func returns  a language name of a given request.
func getLang(r *http.Request) string {
	var lang string
	lang = getOptions(r)["lang"]
	if lang == "" {
		panic(fmt.Sprintf("getLang: The lang key is empty for %s", r.RequestURI))
	}
	return lang
}

//220118 05:32 It creates an account of only type 0
func regAccount(r *http.Request) {
	var accName string = accountName(r)
	if getAccount(accName) != nil {
		return
	} else {
		var opts = make(map[string]string)
		opts["lang"] = "ru"
		var newAcc = &Account{accName, 0, opts, time.Now()}
		accountsMtx.Lock()
		defer accountsMtx.Unlock()
		accounts.PushFront(newAcc)
	}
}

func saveAccounts(accountsFileName string) {
	var f *os.File
	var err error
	var account *Account
	var line string
	if f, err = os.Create(accountsFileName); err != nil {
		panic(fmt.Sprintf("saveAccounts: creating file %s err=%s", accountsFileName, err.Error()))
	}

	accountsMtx.Lock()
	defer accountsMtx.Unlock()

	for e := accounts.Front(); e != nil; e = e.Next() {
		account = e.Value.(*Account)
		line = fmt.Sprintf("%s;%s;at %s\n", account.Name, account.Options, account.RegTm.Format("20060102_150405"))
		f.WriteString(line)
	}
}
