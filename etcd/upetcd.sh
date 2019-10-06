#!/bin/bash
a=`curl  -s "http://ip.tool.chinaz.com/" | grep -1n '您的IP' |grep -v '您的IP'| awk -F ">" '{print $2}' | awk -F "[</dt]" '{print $1}'| awk 'NR==2{print $1}'`
old=`cat /root/command/ip`
echo $a > /root/command/ip
if [ "$a" = "$old" ]
then
        echo "same"
else
        /root/command/upetcd  $a
        curl -X POST http://111.231.82.173:9000/save -H "application/x-www-form-urlencoded" -d "ip=${a}"
fi