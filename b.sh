#!/bin/bash

#220124 18:36 + stopping and starting the server on 95.213.191.152


appname="${PWD##*/}"

echo 220104 10:57 220106 17:16 220124 18:36 It is for 1. building $appname 2. stopping  and copying 3. starting
compiltime=$(date +%y%m%d_%H%M)


sed -i "s/---.*---/---$appname from $compiltime---/" mn.go
#sed -i "s/---.*---/---$appname---/" ind.html
#sed -i "s/:::.*:::/:::$compiltime:::/" ind.html

if [ $? != 0 ]; then 
echo the sed failed
exit
fi

go build
if [ $? != 0 ]; then 
echo the go build failed
exit
fi

#scp -r $appname html image b.sh set_ind.sh *.go favicon.ico mmsit.js mmsite qqmak@192.168.1.44:~/Progects/freelancer/mmsite 
#if [ $? != 0 ]; then 
#echo the copying to 44 failed
#exit
#fi
echo stopping start ----------------------------------
ssh root@95.213.191.152 "pkill mmsite;"
echo stoping end -------------------------------------

echo scp start ==================================

scp -r $appname html image  b.sh set_ind.sh *.go favicon.ico mmsit.js root@95.213.191.152:~/mmsite 
if [ $? != 0 ]; then 
echo the copying to Vscale failed
#exit
fi


echo scp end ==================================


echo deleting files start ----------------------------------
ssh root@95.213.191.152 "cd mmsite;rm *.log nohup.out"
echo deleting files end -------------------------------------

echo launching start ----------------------------------
echo no realization yet ssh root@95.213.191.152 "cd mmsite;nohup ./mmsite & ;pgrep mmsite "
echo launching end -------------------------------------



#rm $appname
echo Ok







