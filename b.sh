#!/bin/bash

#echo 220202 06:29 I do not know yet what it does
echo 220411 08:57 b.sh - the local building of the mmst server with further manipulations on root@95.213.191.152 
echo 

appname=mmst
hostname=$(hostname)
compiltime=$(date +%y%m%d_%H%M)
last_git_commit=$(git log --pretty=format:"%h" -n 1)
git_branch=$(git branch | sed 's/ //g' )
last_git_commit_tag=$(git describe --tags $last_git_commit)
versionInfo=$appname---$last_git_commit_tag---$hostname---$compiltime





go build -ldflags "-X main.versionInfo=$versionInfo" -o $appname
if [ $? != 0 ]; then 
echo 1:golang building failed;
exit;
else
echo 1:The compilation was successful. ver: $versionInfo;
fi




sed -i "s/:::.*:::/:::$versionInfo:::/" html/*.html mmsit.js mystyle.css
if [ $? != 0 ]; then 
echo 2:the sedding failed
exit
else
echo 2:The version was sedded to html/*.html mmsit.js mystyle.css
fi

ssh root@95.213.191.152 "pkill $appname"
if [ $? != 0 ]; then 
echo 3: the stopping the old version of $appname was failed.

else
echo 3: the stopping the old version of $appname was successful.
fi


scp -r $appname html image  b.sh *.go favicon.ico mmsit.js mystyle.css root@95.213.191.152:~/mmst 
if [ $? != 0 ]; then 
echo 4: the scping project files to the cloud was failed.
else
echo 4: the scping project files to the cloud was successful.
fi


ssh root@95.213.191.152 "cd mmst;rm *.log out.txt"
if [ $? != 0 ]; then 
echo 5: the removing old logs and nohup.out was failed.

else
echo 5: the removing old logs and nohup.out was successful.
fi




ssh root@95.213.191.152 "cd mmst; ./mmst mode=11 &"
if [ $? != 0 ]; then 
echo 6: the launching of the new version was failed.
exit
else
echo 6: the launching of the new version was successful.
fi




echo Ok







