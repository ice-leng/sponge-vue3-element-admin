#!/bin/bash

# chkconfig: - 85 15
# description: admin

serverName="admin"
cmdStr="cmd/${serverName}/${serverName}"
pidFile="cmd/${serverName}/${serverName}.pid"
configFile=$1

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

function stopService(){
    local NAME=$1

    # priority to kill process by pid
    if [ -f "${pidFile}" ]; then
        local pid=$(cat "${pidFile}")
        local processInfo=`ps -p "${pid}" | grep "${cmdStr}"`
        if [ -n "${processInfo}" ]; then
           kill -9 ${pid}
           checkResult $?
           echo "Stopped ${NAME} service successfully, process ID=${pid}"
           rm -f ${pidFile}
           return 0
        fi
    fi

    # if the pid file does not exist, get the pid from the process name and kill the process
    ID=`ps -ef | grep "${cmdStr}" | grep -v "$0" | grep -v "grep" | awk '{print $2}'`
    if [ -n "$ID" ]; then
        for id in $ID
        do
           kill -9 $id
           echo "Stopped ${NAME} service successfully, process ID=${ID}"
           return 0
        done
    fi
}

function startService() {
    local NAME=$1

    sleep 0.2
    go build -o ${cmdStr} cmd/${NAME}/main.go
    checkResult $?

    # running server, append log to file
    if test -f "$configFile"; then
        nohup ${cmdStr} -c $configFile >> ${NAME}.log 2>&1 &
    else
        nohup ${cmdStr} >> ${NAME}.log 2>&1 &
    fi

    local pid=$!
    printf "%s" "${pid}" > "${pidFile}"
    sleep 1

    local processInfo=`ps -p "${pid}" | grep "${cmdStr}"`
    if [ -n "${processInfo}" ]; then
        echo "Started the ${NAME} service successfully, process ID=${pid}"
    else
        echo "Failed to start ${NAME} service"
        rm -f ${pidFile}
		    return 1
    fi
    return 0
}

stopService ${serverName}
if [ "$1"x != "stop"x ] ;then
    sleep 1
    startService ${serverName}
    checkResult $?
else
    echo "Service ${serverName} has stopped"
fi
