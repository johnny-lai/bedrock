#!/bin/bash

i="0"

while [ $(kubectl get svc $1 -o json | jq -r '.status.loadBalancer.ingress[0].ip') = "null" ] && [ $i -lt 20 ]; do
	kubectl get pods
	kubectl describe svc $1
	echo [$i] Waiting for $1 to be ready...
	sleep 5;
	i=$[$i+1]
done

if [ $i -eq 20 ]; then
  echo "Timeout while waiting for $1 to be ready"
fi

kubectl describe svc $1
