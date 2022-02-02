#!/bin/bash

#As it is well known that a http(s) system has two components, namely a backend and frontend.
#The mmst is such system.
#Its backend is a https server placed on a cloud virtual host (or on another place accordingly the developer's momentary needs).
#Its another component (on now, 220201 06:26) is a browser page that builds from a set of files.
#This script copies the set from a worker area of the developer to a productive area of the system (now root@95.213.191.152:~/mmst).

#appname="${PWD##*/}"

echo 220124 04:49 220121 05:49 17:16 It is for transferring the frontend files to the 95.213.191.152 virtual host.
echo -----------------------------------

starttime=$(date +%y%m%d_%H%M)
echo Start:$starttime ==================================

scp -r html image  favicon.ico mmsit.js root@95.213.191.152:~/mmst 



endtime=$(date +%y%m%d_%H%M)
echo End:$endtime ==================================










