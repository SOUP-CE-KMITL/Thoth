#!/bin/bash
./setup/kube_start.sh
./setup/influx_start.sh
./setup/haproxy/vamp_start.sh
sleep 10

template=$(<setup/haproxy/vampConfig.template.json)
for i in {1..10};
do
#echo $template
        let Y=$(($i+9000))
        let Z=$((30000+$i))
        res=$(sed "s/X/$i/g" <<< $template)
        res=$(sed "s/Y/$Y/g" <<< $res)
        res=$(sed "s/Z/$Z/g" <<< $res)
        echo $res | http POST localhost:10001/v1/routes
done

