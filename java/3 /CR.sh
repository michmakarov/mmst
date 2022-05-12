#!/bin/bash

project="Globals.java HttpServer.java SocketProcessor.java FrontLog1.java FileReader.java"
compiletime=$(date +%y%m%d_%H%M)

echo Project : $project
echo $compiletime : The project will be compiled and started

javac -Xdiags:verbose $project

if [ $? != 0 ]; then 
echo Compiling failed;
exit;
else
echo The compilation was successful. Next trying to start HttpServer.;
fi

java HttpServer
