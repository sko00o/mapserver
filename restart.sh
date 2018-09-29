#!/bin/bash
app=map
pkill -9 $app
go build -o $app
nohup ./$app -c appserver.conf > $app.log 2>&1 &
ps -eo comm,lstart | grep $app
