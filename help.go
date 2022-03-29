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
<p>For having asscees to the site you have to obtain a cookie.
For this you should give a special request: <a href="/registerme"> registration </a>
</p>
`
var CookieIs = `
<h2> You have been registered</h2>
<p>Now you may do any request
</p>
`
