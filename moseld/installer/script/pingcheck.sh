#!/usr/bin/env bash

NODENAME=$1
URL=$2
PINGMS=

[ -z $NODENAME ] && exit 1 || [ -z $URL ] && exit 1 

ping -c 1 $URL &> /dev/null
if [[ $? != 0 ]]; then
  echo "state:FAIL"
  exit 1; 
fi

PINGMS=$(ping -c 4 $URL | tail -1| awk -F '/' '{print $5}')
echo "state:OK"
echo "time:$PINGMS"
