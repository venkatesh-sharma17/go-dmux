#!/bin/bash

mkdir -p  /var/log/flipkart/go-dmux

#if [[ -z ${LOG_DIR} ]]; then
#  export LOG_DIR="/var/log/go-dmux"
#  mkdir -p ${LOG_DIR}
#fi
sleep 10
/app/go-dmux /app/config/conf.json
