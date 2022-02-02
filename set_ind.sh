#!/bin/bash
#220106 17:37 220201 06:12
#It is invoked at runtime by the server before listening of connections
#exit 123
appname="${PWD##*/}"

compiltime=$(date +%y%m%d_%H%M)
echo set_ind.sh : $appname of $compiltime
sed -i "s/---.*---/---$appname---/" html/ind_ru.html html/ind_en.html
sed -i "s/:::.*:::/:::$compiltime:::/" html/ind_ru.html html/ind_en.html
