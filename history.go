// history
package main

var history = `
220104 11:24<br>
Here new features of the mmsite will be described.<br>
Still here no describing. For getting info about you have only /about
+++++++++++++++++++++++++++++++++++++++++++++++++++++++<br>
220105 13:18 <br> 
Must be permit of multiple args; see func setArgs; mode <br>
It is worth to say that the old version of setArgs in whole fits that requirement.<br>
Only new global variables must be created.In that cause it is the var mode = 0<br>
+++++++++++++++++++++++++++++++++++++++++++++++++++++++<br>
220106 18:11 <br>
modification ind.html at runtime to show that the file is waring the default program name (that is b.sh was not invoked).<br>
For it a new set_ind.sh was made and the b.sh was changed.<br>
+++++++++++++++++++++++++++++++++++++++++++++++++++++++<br>
220107 05:28<br>
How to log incoming requests? only through some feeler. So let it be.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++<br>
220108 07:11 
1. /favicon.ico
2. The template. It needs to insert at least for decency and decorum.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++<br>
220110 15:46 It seems that scripts is better to hold in distinct file as such decision allows to increase code clearness.
In my mind remembers stirring that I recently have worked out a fit script file but where I will can search it?
It is old case: it is seems that there is somewhere a good thing but where?
+++++++++++++++++++++++++++++++++++++++++++++++++++++++<br>
220112 06:16 What is difference? : aboutAuthor_ru .html vs aboutAuthor_ru.html
It took me about two days to see this
+++++++++++++++++++++++++++++++++++++++++++++++++++++++<br>
220113 04:17 Yesterday with the git nothing went out. Let it be so far. But the enigma of "/aboutAuther"! Let it be called "e_2".
_______07:40 With the e_2 there is some nonsense. But there is a need to widen the notion of the mode for perposes of debugging
_______08:27 The button of changing the language does not work under the firefox!
_______09:23 Something is wrong with my codes, but I do not see what.
_______11:03 Only to repeat remains: no bitterness is more bitter than foolishness
And else: an old horse spoiles furrow as well as a young one but tills more worse
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++<br>
220118 04:45 I have been seeing a rude error: a one can change language and it will be influencing another.
Besides this it is not bad to see who have visited the site at some period.
See accounts.go
_______17:55 As if the error has been mended
+++++++++++++++++++++++++++++++++++++++++++++++
220124 04:41
1. cpFront.sh
2. There is some error with changing the language, but I do not catch it yet.
_______16:55 Let's a name of the error will be langErr220124
_______19:25 with the new version of b.sh nothing turned out
+++++++++++++++++++++++++++++++++++++++++++++++
220125 05:54 "Not repeat code" - the golden rule, but it is very hard to follow it!
It is about from where give a debug message - in handlers or in the feeler.
But there is a question! What to do if a specific debug info is needed?
_______14:30 The langErr220124 has vanished as if. Why? The dick knows it.
But stopping and starting through the b.sh has not turned out as earlier. I leave it for future.
_______15:21 The langErr220124 As soom as I had removed the alerts from the function changeLang(), under the Firefox the button stopped to work
++++++++++++++++++++++++++++++++++++++++++++++
220126 06:26 At last I have seemed to understand the langErr220124.
All the same crap! It shows brilliantly as an old one stamping on the same old rake.
See function changeLang() and function getQueryFunc(uri, func)
_______15:38
1.To password up secret requests
++++++++++++++++++++++++++++++++++++++++++++++
220127 06:18
Temporary /sms?text=...
_______08:45 An sms was sent and received
_______14:51 func normSms(sms string) string  func sendSms(sms, to string) (answer string, err error)
_______17:06 Letters were sent and received with a sms; func letterHandler(w http.ResponseWriter, r *http.Request)
+++++++++++++++++++++++++++++++++++++++++++++++
220128 04:42 The letters
There are needs to limit volume of a letter and capacity of the directory of letters.
_______17:19 As if with the letters all was done.
+++++++++++++++++++++++++++++++++++++++++++++++
220131 08:05 From indHandler:
	//220131 04:36 This nonsense has existed long ago and it had not been noticed up to the last friday.
	// This phenomenon is such interesting that it deserve a name: langErr220131
The phenomenon indeed deserves its own name. But a nonsense was first hurried decision about its bases.
At now (_______08:18) it as if vanished after changing call of of http.ServeFile(w, r, fileName__) to of MyServeFile(w, r, fileName__)
(into func aboutAuthorHandler). Why? I do not know.
______13:15 Description of langErr220131:
The /aboutAuthor is now the only request with two langquage variants besides the / itself and only the last  bear the buttom to change  langquage.
The button change language of the / and it is expected that the /aboutAuthor page changes language after a reload but it does not!
______13:41 With Chome the err is not.
______13:45 With Fox the err is not.
______13:50 With Opera the err is not. Furthermore, if the change is done in the Chrome, it has correct outcome in the two others.
++++++++++++++++++++++++++++++++++++++++++++++++
220201 04:39 Now the directory is named mmst
The problem: On the github (https://github.com/michmakarov/mmsite) there is the key file, that is a private key.
It is not seemed well: I cannot make the repository public.
I did not found a better decision than to create a new repository with name "mmst"
`