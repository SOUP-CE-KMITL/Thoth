#!/bin/bash
#docker pull tutum/influxdb:0.10
#kubectl create -f influxDB.yaml 

docker run -d --privileged --restart=always --name=thoth-influxdb -v /var/lib/influxdb:/data -p 8083:8083 -p 8086:8086 -e ADMIN_USER="thoth" -e INFLUXDB_INIT_PWD="thoth" -e PRE_CREATE_DB="thoth;" tutum/influxdb:0.10
