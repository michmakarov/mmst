// 220117 18:26 accounts
package main

import (
	"container/list"
	"fmt"
	"io/ioutil"

	//"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Account struct {
	Name    string            // identificator
	Tp      int               //type 0 - automatically created; the name is remote address; no password
	Options map[string]string //options, sequence of substring of format key=value,for example lang=en somekey=somevalue ...
	//220303 07:30 the key and the value must be a valid utf-8, and must not contain runes that is <= 20
	//220307 16:08 see func strToOpts(s string) (opts map[string]string , optsToStr(opts map[string]string) string , and func goodString(s string) (err error)
	RegTm time.Time //Time of regidtrarion
}

var accounts *list.List
var accountsMtx sync.Mutex // it queues queries for access to the accounts

func init() {
	if accounts == nil {
		accounts = list.New()
	}
	restoreAccounts()
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
//220325 05:48 Now a remote address is not an account name; see var accountName in feeler.go ((f *feeler) ServeHTTP)
//It is like the getOptions function. For the good in both must be used func getAccount2(r *http.Request) that panics if there is not an account
func accountName(r *http.Request) string {
	var accName string
	var accRes byte
	if accName, accRes = getCookieVal(r); accRes != 0 {
		panic(fmt.Sprintf("getOptions (accaouts.go): getting cookie value problem, accRes=%d", accRes))
	}
	if getAccount(accName) == nil {
		panic(fmt.Sprintf("getOptions (accaouts.go): no such account %s", accName))
	}
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	for e := accounts.Front(); e != nil; e = e.Next() {
		if e.Value.(*Account).Name == accName {
			return e.Value.(*Account).Name
		}
	}
	panic(fmt.Sprintf("getCurrOptions: no account for %s", r.RequestURI))
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
//func getOptions(r *http.Request) map[string]string {
//220325 07:57 Now the accName is retrieved from the r through the getCookieVal func (cookie. go)
func getOptions(r *http.Request) map[string]string {
	var accName string
	var accRes byte
	//if accName, _, err = net.SplitHostPort(r.RemoteAddr); err != nil {
	//	panic(fmt.Sprintf("getCurrOptions: net.SplitHostPort err=%s", err.Error()))
	//}
	if accName, accRes = getCookieVal(r); accRes != 0 {
		panic(fmt.Sprintf("getOptions (accaouts.go): getting cookie value problem, accRes=%d", accRes))
	}
	if getAccount(accName) == nil {
		panic(fmt.Sprintf("getOptions (accaouts.go): no such account %s", accName))
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
//220325 07:34 It creates a new account if there is not one with given name aN
//It take a remote address (RA) and an user agent (UA) as options
func regAccount(aN string, r *http.Request) {
	//var accName string = accountName(r)
	if getAccount(aN) != nil {
		return
	} else {
		var opts = make(map[string]string)
		opts["lang"] = "ru"
		opts["RA"] = r.RemoteAddr
		opts["UA"] = r.UserAgent()
		var newAcc = &Account{aN, 0, opts, time.Now()}
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

//220302 17:10
//It saves accout list (accounts) as text file (see global var accountsFileName)
//220308 09:51 Eche line of the file (without \n) represents an account in format:
// <account name>;<account type>;<account options (see utils.optsToStr func)>;<time of registration>
func saveAccountList() {
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
		line = fmt.Sprintf("%s;%d;%s;%s\n", account.Name, account.Tp, optsToStr(account.Options), account.RegTm.Format("20060102_150405"))
		f.WriteString(line)
	}
}

//220307 16:38
//220308 08:56 It is presumed and checked that the file consist of lines that have format that is described for saveAccountList function.
//220309 09:54
//220310 16:46 there was a big fuss about reading file content. See further matter into Progects/golang/220310_rf
func restoreAccounts() {
	var f *os.File
	var fi os.FileInfo
	var err error
	var buf []byte
	var lines []string
	var ac *Account
	var acc int //220310 16:46 counter of restored accounts

	if accounts.Len() != 0 {
		panic("Illegal call of restoreAccounts function: it may be called when the accounts list is empty.")
	}

	if f, err = os.Open(accountsFileName); err != nil {
		panic(fmt.Sprintf("restoreAccounts: opening file %s err=%s", accountsFileName, err.Error()))
	}

	if fi, err = f.Stat(); err != nil {
		panic(fmt.Sprintf("restoreAccounts: getting info of %s; err=%s", accountsFileName, err.Error()))
	}
	if fi.Size() > int64(accountsFileMaxSize) {
		WriteToLog(fmt.Sprintf("Size of %s is more than %d; accounts will be accounted newly", accountsFileName, accountsFileMaxSize))
		return
	}

	//buf = make([]byte, accountsFileMaxSize)

	//if content, err = os.ReadFile(accountsFileName); err != nil {
	//	panic(fmt.Sprintf("restoreAccounts: reading content of file %s err=%s", accountsFileName, err.Error()))
	//}

	//if _, err = f.Read(buf); err != nil {
	//	panic(fmt.Sprintf("restoreAccounts: reading content of file %s err=%s", accountsFileName, err.Error()))
	//}

	//if err = io.ReadFull(f); err != nil {
	//	panic(fmt.Sprintf("restoreAccounts: reading content of file %s err=%s", accountsFileName, err.Error()))
	//}

	if buf, err = ioutil.ReadFile(accountsFileName); err != nil {
		panic(fmt.Sprintf("restoreAccounts: reading content of file %s err=%s", accountsFileName, err.Error()))
	} //else {
	//		fmt.Printf("%s\n----------------\n", string(buf))
	//	}

	lines = strings.Split(string(buf), "\n")

	for _, line := range lines {

		if line == "" {
			continue
		} else {
			acc++
		}
		ac = accountLineToAccount(line)
		accounts.PushFront(ac)
	}
	fmt.Printf("restoreAccounts: %d accounts was successfully restored\n", acc)

}

//220309 10:42
func accountLineToAccount(l string) (ac *Account) {
	var acSl []string
	acSl = strings.Split(l, ";")
	if len(acSl) != 4 {
		panic(fmt.Sprintf("accountLineToAccount: an account line must have 4 conponents, but there is line \"%s\"", l))
	}
	ac = new(Account)
	ac.Name = checkAccontName(acSl[0])
	ac.Tp = convAccontTp(acSl[1])
	ac.Options = strToOpts(acSl[2])
	ac.RegTm = convAccountTm(acSl[3])

	return ac
}

//220325 09:12 The func removes the corresponding account if it exists
func delAccount(accName string) {
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	for e := accounts.Front(); e != nil; e = e.Next() {
		if e.Value.(*Account).Name == accName {
			accounts.Remove(e)
			return
		}
	}
	return
}
