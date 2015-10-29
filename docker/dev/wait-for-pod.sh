#!/bin/bash

i="0"

while [ $(kubectl get pods -l name=$1 -o json | jq "reduce .items[].status.containerStatuses[].ready as \$ready (true; . and \$ready)") != "true" ] && [ $i -lt 10 ]; do
	kubectl get pods;
	echo Waiting for $1 to be ready...
	sleep 3;
	i=$[$i+1]
done

kubectl get pods
