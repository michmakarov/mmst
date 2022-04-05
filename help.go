//220106 16:39 For explanation of the program (its design and properties)
package main

var help = `
This is a CLI program that realizes HTTP or HTTPS server dependly of mode.
The common format of a command is prog agr1 arg2 ... .
The effect of an argument depends on its nature and describes further.
ARGUMENTS:
h - it must be first
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
