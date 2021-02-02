#!/bin/bash

mkdir -p  /var/log/go-dmux
sleep 10
/app/go-dmux /app/config/conf.json > /var/log/go-dmux/stdout.log 2> /var/log/go-dmux/stderr.log
