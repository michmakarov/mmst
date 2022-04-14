//220106 16:39 For explanation of the program (its design and properties)
package main

var help = `
<html>
<head>
<style>
  h1,h2,h3 {margin-bottom: 0px;}
</style>
</head>
<body>
<h1>The http server of the mmst site has next features.</h1>
<h2>VERSIONING</h2>
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
<h2>FRONT LOG</h2>
About each incoming request the server make a record into a special file, that is named "Front log"<br>
The record has format: "property--property--...property" <br>
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
For example <a href="/">Start page</a> will show index page (/) of the site 
</p>
`

var CookieIsStill = `
<h2> You have been registered already</h2>
<p>
You have an account on the site. It has been prolonged for 720 hours.<br>
</p>
`
