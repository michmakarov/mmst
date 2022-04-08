#!/bin/bash

echo 220407 07:43 bl.sh - the local building of mmst
#bl stands for build loccaly. Initially it is copy of b.sh
#The main desire of it was to show clearly (on /ind page) the version of executable with its git data.
#But, as usual, I fogot to save it and started the initial content, that is in fact I started b.sh
# So, now on the cloud works the executable with tag=0.0 and with little intelligible info

appname="${PWD##*/}"
ndname=$(uname -n)
compiltime=$(date +%y%m%d_%H%M)
last_git_commit=$(git log --pretty=format:"%h" -n 1)
git_branch=$(git branch | sed 's/ //g' )
last_git_commit_tag=$(git describe --tags $last_git_commit)
versionInfo=$appname-of-$last_git_commit_tag,built-on-$ndname-at-$compiltime-commit:$last_git_commit


#echo appname=$appname--
#echo ndname=$ndname--
#echo compilation=$compiltime--
#echo last_git_commit=$last_git_commit--
#echo git_branch=$git_branch--
#echo last_git_commit_tag=$last_git_commit_tag--
#echo vi:$versionInfo

#exit




go build -ldflags "-X main.versionInfo=$versionInfo"

if [ $? != 0 ]; then 
echo golang building failed;
exit;
else
echo The compilation was successful;
echo ver: $versionInfo;
fi


sed -i "s/:::.*:::/:::$versionInfo:::/" html/*.html mmsit.js mystyle.css


