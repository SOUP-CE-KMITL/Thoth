#!/bin/bash
./setup/kube_start.sh
./setup/influx_start.sh
./setup/haproxy/vamp_start.sh
http POST localhost:10001/v1/routes < setup/haproxy/vamp_config.json
http POST localhost:10001/v1/routes < setup/haproxy/vamp_config2.json
