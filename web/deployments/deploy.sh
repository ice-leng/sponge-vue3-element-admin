#!/bin/bash

serviceName="admin"

function checkResult() {
    result=$1
    if [ ${result} -ne 0 ]; then
        exit ${result}
    fi
}

# determine if the startup service script run.sh exists
runFile="/app/${serviceName}-binary/run.sh"
if [ ! -f "$runFile" ]; then
  # if it does not exist, copy the entire directory
  mkdir -p /app/${serviceName}-binary
  cp -rf /tmp/dist /app/${serviceName}-binary
  checkResult $?
  rm -rf /tmp/dist*
else
  # replace only the binary file if it exists
  rm -rf /app/${serviceName}-binary/dist
  cp -rf dist /app/${serviceName}-binary/dist
  checkResult $?
  rm -rf /tmp/dist*
fi

rm -rf /app/${serviceName}-binary/dist/deploy.sh

echo "web directory is /app/${serviceName}-binary/dist"
