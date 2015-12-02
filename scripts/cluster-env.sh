#!/bin/bash
CLUSTER_TYPE=`kubectl get svc go-service-basic -o json | jq -r '.spec.type'`

if [ "$CLUSTER_TYPE" = "NodePort" ]; then
  CLUSTER_IP=`kubectl get nodes -o json | jq -r '.items[0].spec.externalID'`
  CLUSTER_PORT=`kubectl get svc go-service-basic -o json | jq -r '.spec.ports[0].nodePort'`
fi

export CLUSTER_TYPE
export CLUSTER_IP
export CLUSTER_PORT

env | grep "^CLUSTER_"