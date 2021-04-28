#!/bin/bash

mkdir -p  /var/log/flipkart/go-dmux

#if [[ -z ${LOG_DIR} ]]; then
#  export LOG_DIR="/var/log/go-dmux"
#  mkdir -p ${LOG_DIR}
#fi
sleep 10
# arguments for the app will be passed at runtime
exec /app/go-dmux "$@"
