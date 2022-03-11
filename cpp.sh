#!/bin/bash

#220302 09:09 This opens "prog.cpp", interpret it as a cpp code, and trying to conpile and run it.


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

g++ cpp/hw.cpp
if [ $? != 0 ]; then 
exit;
else
./a.out
fi







