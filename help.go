//220106 16:39 For explanation of the program (its design and properties)
package main

var help = `
<html>
<head>
<style>
  h1,h2,h3 {margin-bottom: 0px;}
  a {margin-left: 150px;}
</style>
</head>
<body>
<h1>The http server of the mmst site has next features.</h1>
	<a href="#h2.1">Invoking and mode.</a><br>
	<a href="#h2.2">VERSIONING</a><br>
	<a href="#h2.3">FRONT LOG</a><br>
	<a href="#h2.4">GENERAL LOG</a><br>
	<a href="#h2.5">ACCOUNTING</a><br>
	
<h2 id="h2.1">Invoking and mode.</h2>
The server is CLI program that has next command line for starting:<br>
prog mode=&lt;mode value&gt;. Here inside angle brackets  there is a decimal integer that represents a mode;for example ./mmst mode=1<br>
Here the prog stands for the path to a executable file of the server.<br>
The mode is a number in decimal notation, where each decimal place that is numbered from less significal to more ones and first has number 0<br>
The figure in a decimal place codes some feature so that 0 tells absence of the feature.<br>
So, "prog mode=0" starts the server with features that are switched off<br>
There are next features:<br>
<h3>0 - debugging</h3>  It codes rules of printing out a diagnostic messages.<br>
Particularly, if it value > 0, the lines of the front log will be doublets to the console (or to a file that stands for the console, see farther)<br>
<h3>1 - HTTP or HTTPS</h3> If in the place is O then HTTP server will be started, othewise - HTTPS
So "prog mode=10" starts HTTPS sever without debuging.<br>
<h3>2 - sms</h3> sending a noting sms to the author when receiving a letter for him.<br>
In the version 1.0 this feature is disabled: in any case is not attempt to send a sms<br>
<h3>3</h3> redirecting stdout and stderr to a file with name "out.txt".<br>
That is if this is more zero and the feature 0 is estated output will be to file regardless there is a console or not<br>
So "prog mode=1001" starts HTTP server with writing diagnostic messages to the file.<br>


<h2 id="h2.2">VERSIONING</h2>
A version of the server is a string with fields divided by "---" (the three going in succession letters "-")<br>
For example "mmst---0.1-1-g7bfa24b---mich412-A320M-S2H-V2---220411_1704"
Versions are constituted and distributed by the command files of bl.sh and b.sl which a formed the same version.<br>
The distribution means that the server version serves as tag that mark other files of the site, such as .html, .js, or .css
The fields have next sense:
<h3>1.Prog name</h3>
it is the base name of an executable file represented the server.<br>
For example "mmst"
<h3>2.TAG</h3>
It is a tag of the last git commit. For example "0.1-1-g7bfa24b"<br>
it is a result of the git command "git describe --tags <last_git_commit"> where <last_git_commit"> is result of "git log --pretty=format:"%h" -n 1"<br>
<h3>3.HOSTNAME</h3>
It is a result of the linux command "hastname", for example mich412-A320M-S2H-V2
<h3>4.TIME</h3>
It represent the buiding time and is result of the linux command "date +%y%m%d_%H%M", for example 220411_1704


<h2 id="h2.3">FRONT LOG</h2>
About each incoming request the server make a record into a special file, that is named "Feeler&lt;20220427_131902&gt;.log"<br>
In angle brackets is a timestamp of creating the file" <br>
Where a property is a pair of name and value divided by "=", as in "ACC=[227 168 102 177 28 184 86 117]"<br>
There are properties:
<h3>DATE</h3>
Its value has format "20220407_160304" (Year,month,day,hour,minute,second) and represent time of obtaining the request<br>
<h3 style="margin-bottom: 0px;">NUM</h3>
It is a order number of the request, for example "NUM=123"<br>
That is there is a counter of incoming requests and its value is the value of the property.<br>
<h3>ACC</h3>
It is a account of the requst.<br
The account is extracted from a cookie named "mmstSession" by decrypting its value<br>
If the process of decrypting was successful the account is in its representation as an array, for example "ACC=[235 206 116 40 240 44 5 224]"<br>
In contra case the property shows error namber of the decrypting process as in "ACC=accRes==2"<br>
There are errors:<br>
1 - no a cookie with name of "mmstSession"<br>
2 - there is cookie but decrypting was not successful<br>
<h3>URI</h3>
For example URI=/accounts?pw=none<br>
<h3>RA</h3>
Remote address. For example RA=[::1]:36174<br>

<h2 id="h2.4">GENERAL LOG</h2>
Into file with name "general.log" the server writes messages of common character, that that the author found interesting for logging<br>
The record has format: "&lt;timestamp&gt;---&lt;file,which generate a record&gt;:&lt;message&gt;" <br>
It is needed to have in mind, that the message may contain the character ":"<br>

<h2 id="h2.5">ACCOUNTING</h2>
For having not trivial access to the site the client must take the cookie through the special request <a href="/registerme">/registerme</a><br>
That request creating an account record and sends a cipher cookie with name of the account, that becomes a name of client's seccion<br>
An account record have not any privat data and using for keeping specific client's options. <br>
<br>


</body>
</html>
`
var noCookieMess = `
<h2> No registration error</h2>
<p>For having access to the site you have to obtain a cookie.
For this you should give a special request: <a href="/registerme"> registration </a>
</p>
`
var CookieIs = `
<h2> You have been registered</h2>
<p>
A cookie with name of "mmstSession" was sent to your agent.
Now you have an account on the site. It will be exist 720 hours.<br>
Your registration code has been sent to you as a cookie with name mmstSession <br>
The cookie is termless but the account may be expired. In this case you must register again<br>
Now you may do any request provide sending with it the cookie and while the account exist
For example <a href="/main">Start page</a> will show main page (/main) of the site 
</p>
`

var youHasAccAlready = `
<h2> You have been registered already</h2>
<p>
You have an account on the site. It has been prolonged for 720 hours.<br>
</p>
`

var yourIPHasAcc = `
<h2> There is account with your IP.</h2>
<p>
Your request for registration has been ignored,
</p>
`
