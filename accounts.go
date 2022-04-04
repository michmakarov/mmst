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
	Name []byte //string 220404 11:14 // identificator
	Tp   int    //type 0 - automatically created; the name is remote address; no password
	//type 1 - since 220330 11:32
	Options map[string]string //options, sequence of substring of format key=value,for example lang=en somekey=somevalue ...
	//220303 07:30 the key and the value must be a valid utf-8, and must not contain runes that is <= 20
	//220307 16:08 see func strToOpts(s string) (opts map[string]string , optsToStr(opts map[string]string) string , and func goodString(s string) (err error)
	RegTm time.Time //Time of regidtrarion
}

var accounts *list.List
var accountsMtx sync.Mutex                      // it queues queries for access to the accounts
var accountTerm, _ = time.ParseDuration("720h") //~ 1 month
//220404 07:52 Separators for forming a text representation of the accounts
var accountCompsSpr = "-;;-"
var optionsSpr = ";--;"
var optionKVspr = ";==;"
var minOptions = 3

func init() {
	if accounts == nil {
		accounts = list.New()
	}
	restoreAccounts()
}

//220118 05:43 The func takes an account name and returns a pointer to an account if it exists or nil if not
func getAccount(accName []byte) *Account {
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	for e := accounts.Front(); e != nil; e = e.Next() {
		if string(e.Value.(*Account).Name) == string(accName) {
			return e.Value.(*Account)
		}
	}
	return nil
}

//220330 12:52
//220330 15:02 It expands the functionallity of the getCookieVal
func getAccount2(r *http.Request) (mess string, res byte) {
	var acc *Account
	mess, res = getCookieVal(r)
	if res != 0 {
		return
	}
	acc = getAccount([]byte(mess))
	if acc == nil {
		res = 3
		mess = fmt.Sprintf("getAccount: from %s was a valid cookie (val = %s) but an account absences", r.RemoteAddr, mess)
	}
	if isDebug(serverMode) {
		fmt.Printf("DDDD---- getAccount2: res=%d; accName=%s\n", res, mess)
	}
	return
}

///* remuved 220331 12:41
//220118 17:12
//220325 05:48 Now a remote address is not an account name; see var accountName in feeler.go ((f *feeler) ServeHTTP)
//It is like the getOptions function. For the good in both must be used func getAccount2(r *http.Request) that panics if there is not an account
//220328 09:37 For what is this function at all?
//_______09:41 It is used in many spots. So there is a need th improve it.
//Let it return "?" if there is any problem with getting the account name from the cookie
//and "???" if from the cookie was obtained a valid name but it is not in the registration list.
//220401 10:49 It extracts the account name from parameter. It panics if the namr is not
func accountName(r *http.Request) string {
	var accName string
	accName = (r.Context().Value(AccNameCtxKey)).(string)
	if accName == "" {
		panic(fmt.Sprintf("getOptions (accaouts.go):  the request gives an empty accName"))
	}
	if accName == "?" {
		panic(fmt.Sprintf("getOptions (accaouts.go):  the request gives accName==?"))
	}
	printDebug(fmt.Sprintf("accountName: URI=%s;accName=%v", r.RequestURI, []byte(accName)))
	return accName
}

//220331 08:49
//220401 10:55 refused
//func accountName2(accName string, accRes int) string {
//	if accRes == 0 {
//		return accName
//	} else {
//		return "?"
//	}
//}

//220118 16:22 The func panics if opts==nil
func setOptions(accName []byte, opts map[string]string) {
	if opts == nil {
		panic(fmt.Sprintf(": no account no opt (nil) for %s", accName))
	}
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	for e := accounts.Front(); e != nil; e = e.Next() {
		if string(e.Value.(*Account).Name) == string(accName) {
			e.Value.(*Account).Options = opts
		}
	}

}

//220118 15:04 The func returns options related a given request.
//If the request is not correspond some account the func panics.
//func getOptions(r *http.Request) map[string]string {
//220325 07:57 Now the accName is retrieved from the r through the getCookieVal func (cookie. go)
//220328 10:27 instead panicing it now return nil
//220331 12:48 almost refuse
//220401 10:23 the r parameter carries now the account name
func getOptions(r *http.Request) map[string]string {
	var accName string
	accName = accountName(r)
	return getOptionsByAcc([]byte(accName))
}

func getOptionsByAcc(accName []byte) (optsCopy map[string]string) {
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	optsCopy = make(map[string]string)
	for e := accounts.Front(); e != nil; e = e.Next() {
		if string(e.Value.(*Account).Name) == string(accName) {
			for key, val := range e.Value.(*Account).Options {
				optsCopy[key] = val
			}
		}
		return
	}
	//panic(fmt.Sprintf("getCurrOptions: no account for %s", r.RequestURI))
	return

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
//220331 12:32 if there is a new account the accounts are saved in file
func regAccount(aN []byte, r *http.Request) {
	//var accName string = accountName(r)
	if getAccount(aN) != nil {
		return
	} else {
		var opts = make(map[string]string)
		opts["lang"] = "ru"
		opts["RA"] = r.RemoteAddr
		opts["UA"] = r.UserAgent()
		opts["HOST"] = r.Host
		var newAcc = &Account{aN, 1, opts, time.Now()}
		accountsMtx.Lock()
		accounts.PushFront(newAcc)
		WriteToLog(fmt.Sprintf("Accont (new) name=%v", []byte(aN)))
		accountsMtx.Unlock()
		saveAccountList()
	}
}

/* Removed 220331 12:35
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
*/

//220302 17:10
//It saves accout list (accounts) as text file (see global var accountsFileName)
//220308 09:51 Eche line of the file (without \n) represents an account in format:
// <account name>;<account type>;<account options (see utils.optsToStr func)>;<time of registration>
//220404 09:05 see var accountCompsSpr(="-;;-") var optionsSpr=(";--;")
//has four substring separeted by the accountCompsSpr
func saveAccountList() {
	var f *os.File
	var err error
	var account *Account
	var line string

	var aCS = accountCompsSpr

	if f, err = os.Create(accountsFileName); err != nil {
		panic(fmt.Sprintf("saveAccounts: creating file %s err=%s", accountsFileName, err.Error()))
	}

	accountsMtx.Lock()
	defer accountsMtx.Unlock()

	for e := accounts.Front(); e != nil; e = e.Next() {
		account = e.Value.(*Account)
		//line = fmt.Sprintf("%s;%d;%s;%s\n", account.Name, account.Tp, optsToStr(account.Options), account.RegTm.Format("20060102_150405"))
		//				   Nm  Tp  op  Tm
		line = fmt.Sprintf("%s%s%d%s%s%s%s\n", byteSliceToStrRepresentation(account.Name), aCS, account.Tp, aCS, optsToStr(account.Options), aCS, account.RegTm.Format("20060102_150405"))
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
		//panic(fmt.Sprintf("restoreAccounts: opening file %s err=%s", accountsFileName, err.Error()))
		//220401 14:09
		fmt.Printf("restoreAccounts: no file %s\n", accountsFileName)
		return
	}

	if fi, err = f.Stat(); err != nil {
		panic(fmt.Sprintf("restoreAccounts: getting info of %s; err=%s", accountsFileName, err.Error()))
	}
	if fi.Size() > int64(accountsFileMaxSize) {
		WriteToLog(fmt.Sprintf("Size of %s is more than %d; accounts will be accounted newly", accountsFileName, accountsFileMaxSize))
		return
	}

	if buf, err = ioutil.ReadFile(accountsFileName); err != nil {
		panic(fmt.Sprintf("restoreAccounts: reading content of file %s err=%s", accountsFileName, err.Error()))
	}
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
//220404 11:32 for what is it?: for restoreAccounts
func accountLineToAccount(l string) (ac *Account) {
	var acSl []string
	acSl = strings.Split(l, accountCompsSpr) // ";")
	if len(acSl) != 4 {
		panic(fmt.Sprintf("accountLineToAccount: an account line must have 4 conponents, but there is line \"%s\"", l))
	}
	ac = new(Account)
	ac.Name = byteStrRepresentationToByteSlice(checkAccontName(acSl[0]))
	ac.Tp = convAccontTp(acSl[1])
	ac.Options = strToOpts(acSl[2])
	ac.RegTm = convAccountTm(acSl[3])

	return ac
}

//220325 09:12 The func removes the corresponding account if it exists
//220330 09:36 It returns true if some accound was deleted
func delAccount(accName string) bool {
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	for e := accounts.Front(); e != nil; e = e.Next() {
		if string(e.Value.(*Account).Name) == string(accName) {
			accounts.Remove(e)
			return true
		}
	}
	return false
}

//220330 10:34
func delExpiredAccounts() (accs []string) {
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	for e := accounts.Front(); e != nil; e = e.Next() {
		if time.Since(e.Value.(*Account).RegTm) >= accountTerm {
			accounts.Remove(e)
			accs = append(accs, byteSliceToStrRepresentation(e.Value.(*Account).Name))
		}
	}
	if isDebug(serverMode) {
		printDebug(fmt.Sprint("delExpiredAccounts:%v", accs))
	}
	return
}

//220330 15:43 It sets the date registration to now if it found an account
func prolongeAccount(accName string) {
	accountsMtx.Lock()
	defer accountsMtx.Unlock()
	for e := accounts.Front(); e != nil; e = e.Next() {
		if string(e.Value.(*Account).Name) == string(accName) {
			e.Value.(*Account).RegTm = time.Now()
		}
	}

}
