#!/bin/bash

while true
do
	echo "Collect"
	sleep 2
	http http://localhost:8182/app/eight-puzzle/metrics >> eight.log
done

