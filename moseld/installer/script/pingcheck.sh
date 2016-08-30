#!/usr/bin/env bash

HOST=$1
PINGMS=

[ -z $HOST ] && exit 1

ping -c 1 ${HOST} &> /dev/null
if [[ $? != 0 ]]; then
  echo "state:FAIL"
  exit 1; 
fi

PINGMS=$(ping -c 4 ${HOST} | tail -1| awk -F '/' '{print $5}')
echo "state:OK"
echo "time:$PINGMS"
