#!/bin/sh


remote_ip=47.106.136.96

#services="wallet_service" 
#services="ws_service" 
#services="user_service" 
services="gateway" 
#services="currency_service user_service gateway"
#services="currency_service user_service price_service gateway"


remote_path="/root/go/src/dig/"


function build_service(){
    for service in $services;
    do
        cd ../$service
        GOOS=linux GOARCH=amd64 go build
    done
}



function push_to_remote(){
    for service in $services;
    do
        scp ../$service/$service root@$remote_ip:$remote_path$service.nw
        result=`ssh root@$remote_ip "cd $remote_path && ls $service.nw"`
        echo "push result:$result"
    done
}

function restart_service(){
    for service in $services;
    do
        service_id=`ps -ef |grep -v grep | grep $service | awk -F " "  '{print $2}'`
        result=`ssh root@$remote_ip "ps -ef | grep -v grep | grep $service"`
        service_id=`echo $result | awk -F " " '{print $2}'`
        echo  "$service $service_id"
        if [ "$service_id" != "" ];then
            result=`ssh root@$remote_ip  "echo 'cd $remote_path && mv -f $service.nw $service && kill -9  $service_id &&  nohup ./$service  > $service.log 2>&1 &' > start_$service.sh && sh start_$service.sh >/dev/null 2>&1 & exit"`
        else
            result=`ssh root@$remote_ip  "echo 'cd $remote_path && mv -f $service.nw $service && nohup ./$service > $service.log  2>&1 & ' > start_$service.sh &&  sh start_$service.sh >/dev/null 2>&1 & exit"`
        fi
        result=`ssh root@$remote_ip " mv -f  start_$service.sh /tmp/" `        
        start_service_result=`ssh root@$remote_ip -f -n "ps -ef | grep -v grep | grep start_$service.sh"`
        start_service_id=`echo $start_service_result| awk -F " " '{print $2}'`
        echo $start_service_id
        if [ "$start_service_id" != "" ];then
            result=`ssh root@$remote_ip "pwd && kill -9 $start_service_id"`
        fi 
        restart_result=`ssh root@$remote_ip "ps -ef | grep -v grep | grep $service"`
        service_id=`echo $restart_result | awk -F " " '{print $2}'`
        echo  "new $service $service_id"        
    done
}


function del_local_service(){
    for service in $services;
    do
        mv -f ../$service/$service /tmp/$service 
    done
}

build_service
push_to_remote
restart_service
del_local_service
