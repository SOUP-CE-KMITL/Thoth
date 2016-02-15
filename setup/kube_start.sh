#!/bin/bash
# http://kubernetes.io/v1.1/docs/getting-started-guides/docker.html
#docker run --net=host -d gcr.io/google_containers/etcd:2.0.12 /usr/local/bin/etcd --addr=127.0.0.1:4001 --bind-addr=0.0.0.0:4001 --data-dir=/var/etcd/data
#docker run --net=host --privileged -d -v /sys:/sys:ro -v /var/run/docker.sock:/var/run/docker.sock  gcr.io/google_containers/hyperkube:v1.0.1 /hyperkube kubelet --api-servers=http://localhost:8080 --v=2 --address=0.0.0.0 --enable-server --hostname-override=127.0.0.1 --config=/etc/kubernetes/manifests --read-only-port=10255
#docker run -d --net=host --privileged gcr.io/google_containers/hyperkube:v1.0.1 /hyperkube proxy --master=http://127.0.0.1:8080 --v=2

#docker run -d -p 8083:8083 -p 8086:8086 --expose 8090 --expose 8099 --name influxsrv tutum/influxdb
#-----------------------------------------
#wget https://s3.amazonaws.com/influxdb/influxdb-0.10.0-1.x86_64.rpm
#sudo yum localinstall influxdb-0.10.0-1.x86_64.rpm

##--------------------
docker run --net=host -d gcr.io/google_containers/etcd:2.2.1 /usr/local/bin/etcd --addr=127.0.0.1:4001 --bind-addr=0.0.0.0:4001 --data-dir=/var/etcd/data
#docker run \
#    --volume=/:/rootfs:ro \
#    --volume=/sys:/sys:ro \
#    --volume=/dev:/dev \
#    --volume=/var/lib/docker/:/var/lib/docker:ro \
#    --volume=/var/lib/kubelet/:/var/lib/kubelet:rw \
#    --volume=/var/run:/var/run:rw \
#    --net=host \
#    --pid=host \
#    --privileged=true \
#    -d \
#    gcr.io/google_containers/hyperkube:v1.0.1 \
#    /hyperkube kubelet --containerized --hostname-override="127.0.0.1" --address="0.0.0.0" --api-servers=http://localhost:8080 --config=/etc/kubernetes/manifests

# v1.1.7
#gcr.io/google_containers/hyperkube:v${K8S_VERSION}
#gcr.io/google_containers/hyperkube-amd64:v${K8S_VERSION}
export K8S_VERSION=1.1.7
export MASTER_IP=10.0.1.17
export ETCD_VERSION=2.2.1
export FLANNEL_VERSION=0.5.5
export FLANNEL_IFACE=bridge0
export FLANNEL_IPMASQ=true

# https://github.com/kubernetes/kubernetes/blob/release-1.1/docs/getting-started-guides/docker-multinode/master.md
docker run \
    --volume=/:/rootfs:ro \
    --volume=/sys:/sys:ro \
    --volume=/dev:/dev \
    --volume=/var/lib/docker/:/var/lib/docker:rw \
    --volume=/var/lib/kubelet/:/var/lib/kubelet:rw \
    --volume=/var/run:/var/run:rw \
    --net=host \
    --privileged=true \
    --pid=host \
    -d gcr.io/google_containers/hyperkube:v1.1.7 /hyperkube kubelet --api-servers=http://localhost:8080 --v=2 --address=0.0.0.0 --enable-server --hostname-override=127.0.0.1 --config=/etc/kubernetes/manifests-multi --cluster-dns=10.0.0.10 --cluster-domain=cluster.local
docker run -d --net=host --privileged gcr.io/google_containers/hyperkube:v1.1.7 /hyperkube proxy --master=http://127.0.0.1:8080 --v=2
