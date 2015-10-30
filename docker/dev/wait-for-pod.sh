#!/bin/bash

i="0"

while [ $(kubectl get pods -l name=$1 -o json | jq "reduce .items[].status.containerStatuses[].ready as \$ready (true; . and \$ready)") != "true" ] && [ $i -lt 20 ]; do
	kubectl get pods;
	echo Waiting for $1 to be ready...
	sleep 3;
	i=$[$i+1]
done

if [ $i -eq 20 ]; then
  echo "Timeout while waiting for $1 to be ready"
fi

kubectl get pods
