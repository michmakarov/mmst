#!/bin/bash

echo 220202 06:29 I do not know yet what it does


appname="${PWD##*/}"
compiltime=$(date +%y%m%d_%H%M)
last_git_commit=$(git log --pretty=format:"%h" -n 1)
branch=$(git branch | sed 's/ //g' )
last_git_commit_tag=$(git describe --tags $last_git_commit)

if [ $branch != *main ]; then 
echo The branch is not main : $last_git_commit_tag
exit
fi

versionInfo=$(echo $appname_$last_git_commit_tag[branch_$branch,commit_$last_git_commit]_$compiltime)



go build -ldflags "-X main.versionInfo=$versionInfo"

if [ $? != 0 ]; then 
echo golang building failed;
exit;
else
echo The compilation was successful;
echo ver: $versionInfo;
fi


sed -i "s/:::.*:::/:::$versionInfo:::/" html/*.html mmsit.js mystyle.css

#exit



if [ $? != 0 ]; then 
echo the sed failed
exit
fi

echo stopping mmst start ----------------------------------
ssh root@95.213.191.152 "pkill mmst;"
echo stoping mmst end -------------------------------------


echo scp start ==================================
scp -r $appname html image  b.sh *.go favicon.ico mmsit.js mystyle.css root@95.213.191.152:~/mmst 
echo scp end ==================================


echo deleting files start ----------------------------------
ssh root@95.213.191.152 "cd mmst;rm *.log nohup.out"
echo deleting files end -------------------------------------

echo launching start ====================================
echo no realization yet ssh root@95.213.191.152 "cd mmsite;nohup ./mmsite & ;pgrep mmsite "
echo launching end ======================================




echo Ok







