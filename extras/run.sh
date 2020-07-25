#!/bin/bash

export HOSTIP=$(hostname -I | awk '{print $1}')
echo "Setting env variable $HOSTIP"

docker-compose -f extras/docker-compose.yaml up