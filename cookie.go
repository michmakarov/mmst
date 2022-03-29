// cookie
// 220322 04:27
// 220325 04:36 after /Progects/golang/220315_crypto/5-AES
package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"net/http"

	crand "crypto/rand"
	"crypto/sha256"
	"hash"

	"encoding/hex"
	"fmt"
	"math/rand"
	//"os"
	//"time"
)

const cookieName = "mmstSession"

//------ Global variables that are established by the init function
var key = []byte("12345678901234567890123456789012")
var block cipher.Block
var blockSize int
var iv []byte
var mac hash.Hash

//---------------------

func init() {
	var err error
	if block, err = aes.NewCipher(key); err != nil {
		panic(fmt.Sprintf("aes.NewCipher(key) err=%s", err.Error()))
	}
	blockSize = block.BlockSize()
	iv = make([]byte, blockSize)
	if _, err := rand.Read(iv); err != nil {
		panic(fmt.Sprintf("init (cookie.go):rand.Reader(iv) err=%s", err.Error()))
	}
	mac = hmac.New(sha256.New, key)
	//fmt.Printf("------------------cookie.go init macsize=%d; blockSize=%d", mac.Size(), blocksize)
}

//220329 04:11 before returning it turns the binary text to alphabetical one
func encrypt(plainText []byte) (cipherText []byte) {
	var abCipherText []byte // the alphabetical form of the cipherText
	var plainTextMAC []byte

	mac.Reset()
	mac.Write(plainText)
	plainTextMAC = mac.Sum(nil)

	stream := cipher.NewCFBEncrypter(block, iv)
	plainText = append(plainText, plainTextMAC...)
	cipherText = plainText
	stream.XORKeyStream(cipherText, plainText)

	abCipherText = make([]byte, hex.EncodedLen(len(cipherText)))
	hex.Encode(abCipherText, cipherText)
	cipherText = abCipherText
	return
}

//220329 04:25 as first step it decodes the given text from alphabetical form to binary one
func decrypt(abCipherText []byte) (pT []byte, err error) {
	var cipherText []byte //in binary form
	var MAC1, MAC2 []byte
	var indMAC int // The index where first byte of the MAC is in decrypted cipherText

	cipherText = make([]byte, hex.DecodedLen(len(abCipherText)))
	if _, err = hex.Decode(cipherText, abCipherText); err != nil {
		pT = nil
		err = fmt.Errorf("decrypt(cookie.go): decoding into binary err=%s", err.Error())
		return

	}
	stream := cipher.NewCFBDecrypter(block, iv)

	if (len(cipherText) - mac.Size()) < 0 {
		pT = nil
		err = fmt.Errorf("decrypt(cookie.go): bad length of cipherText(%d)(mac.Size=%d)", len(cipherText), mac.Size())
		return
	}

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)
	indMAC = len(cipherText) - mac.Size()
	if isDebug(serverMode) {
		fmt.Printf("-----------------------decrypt: indMAC=%d; len(cipherText)=%d;  mac.Size()=%d\n", indMAC, len(cipherText), mac.Size())
	}
	MAC1 = cipherText[indMAC:]
	pT = cipherText[:indMAC]
	mac.Reset()
	mac.Write(pT)
	MAC2 = mac.Sum(nil)
	if !hmac.Equal(MAC1, MAC2) {
		pT = nil
		err = fmt.Errorf("decrypt(cookie.go): Authentication failed")
	}
	return
}

//220322 16:08
//220325 05:20
// if res==0 it returns the dectrypted cookie value otherwise a error message
//res==0 - the cookie value
//res==1 - No cookie with name==cookieName
//res==2 - problem with decrypting: Authentication failed (but the cookie with given name is)
func getCookieVal(r *http.Request) (mess string, res byte) {
	var cookie *http.Cookie
	var err error
	var buff []byte
	if cookie, err = r.Cookie(cookieName); err != nil {
		res = 1
		mess = err.Error()
		return
	}

	if buff, err = decrypt([]byte(cookie.Value)); err != nil {
		res = 2
		mess = err.Error()
		return
	}
	mess = string(buff)
	return
}

//220323 03:45
//220325 09:41
// It  renerates a account name, encrypts it as a cookie value and sets the cookie for transferring to a client
func setCookie(w http.ResponseWriter) (accountName string) {
	var err error
	var buff []byte = make([]byte, 8)
	var cookieVal []byte
	var cookie http.Cookie
	if _, err = crand.Read(buff); err != nil {
		panic(fmt.Sprintf("setCookie: crand.Read(buff) err=%s", err.Error()))
	}
	cookieVal = encrypt(buff)

	cookie.Name = cookieName
	cookie.Value = string(cookieVal)
	cookie.Path = "/"

	http.SetCookie(w, &cookie)

	accountName = string(buff)
	return
}
