#!/bin/bash

mkdir -p  /var/log/go-dmux
sleep 10
# arguments for the app will be passed at runtime
exec /app/go-dmux "$@"
