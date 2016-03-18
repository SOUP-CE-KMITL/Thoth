#!/bin/bash
docker run -d -v /var/lib/grafana --name grafana-storage busybox:latest
docker run -d -p 3000:3000 --name=thoth-grafana  --volumes-from grafana-storage  grafana/grafana
