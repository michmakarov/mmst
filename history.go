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
______06:16 See "creting_this_rep.txt" for some details.
______08:11 What are the frontend and the backend? There is a urge want to get a clear answer and to make accordingly the versions showing. 
______15:56 As the server (that is the backend) yields all let's say that namely it embodies the system in whole.
Hence the server must generate the version (whatever it means) and spread it to all components of the system.
++++++++++++++++++++++++++++++++++++++++++++++++
220202 04:50 What is system? Are system and application synonyms? In answering we fall into dim pit of utter ambiguity.
So the b.sh generates a version and spreads it to *.html, mmsit.js and mmsit.css
++++++++++++++++++++++++++++++++++++++++++++++++
220204 05:02 What is knowing of a language?
Yesterday I attempted to public the site on https://freelancehunt.com and had a problem.
For desision of it I took "<meta name='freelancehunt' content='cc9a588ee5dcfce'>"
In the context of it the question is emerged - What is a meta tag?
++++++++++++++++++++++++++++++++++++++++++++++++
220302 16:05 There is a puzzle. Let its name be "accoutPuzzle" and it sounds as "How  the accounts list is formed at starting server".
The answer: no way, at starting the list is empty!
The next question: it is good or bad?
If we say that a user should not notice that a new version of the server was started then it is bad.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220307 17:57 The further into forest, the more partisans
Now a need of a log has arisen. Indeed, what to do if file of accounts is too big and all accounts are refused?
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220307 16:26 A new hitch: version of the golang. It is obvious that a new version may bear new functionality!
So the thing is needed to be been controlling.
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220309 09:44
1. About version. Here is 1.15.5. So what? Under another version will be another binary. This is an advantage of compiled languages.
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220311 09:10
1. the current accounting is unfit: the list of accounts is big and therefore hard to control. Why? For example the NAT.
It as it is seemed must give many fictitious accounts. So identifying an account by IP address is a bad idea. But what identify by?
It is come up that the contragent need be tagged by some tag, a cookie is very fit.
2. Why is so big the nohup.out? I think that errors of the system server need out to distinct log file.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220322 05:53 The accounting will be improved. Let's this decision have its own name : 220322-account
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220323 04:16 As if I decide to resolve the 220322-account task too early:
I am not sure that a nance may be the same for all clients (https://pkg.go.dev/crypto/cipher#AEAD)
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220325 07:07 All manipulations with accounts should be in the func (f *feeler) ServeHTTP
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220330 08:44 About the 220322-account
1. A client may or may not support cookies. This needs to be checked
2. I firmly think that the site is not very important. So the "math/rand" is enough for (with a pecular seed) generating random byte sequences
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220331 08:14 myAccountHandler
As if to do decryption for any want of a account name is not decently. So the name must be passed through a context.
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220401 11:11 1. removing old logs
2. About password: I firmly think that if some  malperson has access to the cloud that person may do all if he would want it
_______13:29 The very old question: Is it need to ensure successive access if a process only reading the data?
let's say that to overfart is better than underfart
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220405 08:51 As if all work well.
But restoring accounts appears unnecessary as it demands to store encryption data that does not seem a good idea.
What is concerning the password, let it goes in code still as a variable that maybe in future will be assigned dynamically.
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220407 10:51 What is next?
1. Old question: How to ensure that a client (or, precisely, an agent) is using the server's software (is, css, ...)
2. Request(s) showing log(s)
_______11:24
The first appears now senseless. Maybe in future it will arise again when the question "what is application" will take more clear answer.
The second is more real so let's go doing it
_______14:47 The fact: there is writing into Stdout but there is not into the feeler log
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220414 05:08 After having cleaned up different foolishness I am going to make 1.0
But before it there is need to bring to oder b.sh, namely to improve starting a new version on the cloud. Now. with the nohup, that does not work at all
_______05:42 How to redirect std.out and std.err to a file (as the terminal is absent)
_______11:35 the b.sh does not work fitly: it does not say about remote launching (the six point) although a new version works
that the bash is a very intricate I know. It constantly surprises me and that seems without end
Of cause here to pecularities of the bash the ones of the ssh was added.
There is need to search a new approach but I am tired from the mmst. So a while I leave it.
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220415 08:42 the music did not play for long: feeler panic getLang: The lang key is empty for / (/ from 77.108.87.132:40298)
The problem is reproducing on the Chrome but not on the Opera (if not paying attention to strange behavior of the Opera on 41)
_______11:44 Let's it be problem_220415
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220416 08:11 problem_220415
This has place on both, 44 and 41. I cannot say yet the exact way to repeat it, but it has place exactly.
What to do? More of all afraids me that it has place with Chromium, but not with Opera and Firefox
_______08:35 In first queue there is want of a independent client
_______08:57 But already there is such! see Progects/freelancer/mmclient
Let's copy it under the git of the mmst
_______09:08 Now it is Progects/freelancer/mmst/mmclient
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220419 13:19 the mmclient works
The problem_220415 is with the mmclient. Opera, Chromium but not with FireFox
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220420 09:00 problem_220415
On 41 (42): This is with the Chromium but is not with the FireFox
On 44 the problem is not at all. The last is the most scaring.
The plan: To play with the mmclient on both machines
_______10:48 On 44 under the FireFox a new registration was demanded. Now there are four registration records - two for 44 and two for 42 (former 41).
_______18:58 The fact is that there is doublicating registration records. It is wrong. Why is it occurring?
Let's serialize the feeler by a mutex regardless that I do not see sense in it
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220422 06:20 problem_220415 For more successful catching it is need to do some actions
_______13:38 Absence of the </div><!-- minediv --> tag was foud out in index files.
But it should not cause the problem in any way.
_______14:19 To say more precisely the last was revealed in copy of a unknown now index file namely in the problem_220415.html
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220425 13:10
1. The /registerme is an ordinary request and must be worked out in common way
2. So a mechanism of banning repetition of not performed request must be worked out instead of registration of the /registerme.
+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220426 09:10
1. What is the "banning repetition of not performed request"?
As if in dealing strictly and only with the http protocol that "repetition" simply impossible.
It takes sense if two requests are "the same" in a context that is laying out of the http protocol.
For example, if there is a working session, that exists independently of requests.
In our case there is such thing, the account.
Thus the source address is essential attribute of a http request but it is not so for the account.
We may construct (by many way, I think) and (consequently) speak about many http request belonging to one account.
So the "banning repetition of not performed request" has sence and all ok
2. How to know up an underlying net connection of a given http request?
As if only through the context of a request.
So, https://stackoverflow.com/questions/29531993/accessing-the-underlying-socket-of-a-net-http-response
++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
220428 11:11
1. About an underlying net connection of a given http request. See previous, 2.
	As if the conclusion is wrong. The http.Server takes the marvellous callback function ConnState
2. What is a repetitive call of "/registerme"? If some agent has not provided the (valid) cookie he may has done it deliberately.
	I think the moment must be checked and one or the other decision must be taken.
_______13:53 So it is turned that a tag (that is accounting at all) does not preserve from multiple registration of the same agent
if it can and want to forge its IP. And the NAT here does not do at all.
_______14:17 Although it is very interesting how the NAT deals with the source IP of IP packets
220503 07:07 The question: For web crawler is it matter what protocol is used, with the TSL or without? Answer: no if it is a contemporary crawler.
_______07:31 The actual task is to open needed pages for web crawlers.
******************************************************************************
220505 07:46 Ha! What is difference between a browser and my mmclient?

`
