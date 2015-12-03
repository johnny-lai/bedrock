#!/bin/bash
set_cluster_env () {
  CLUSTER_TYPE=`kubectl get svc $1 -o json | jq -r '.spec.type'`

  if [ "$CLUSTER_TYPE" = "NodePort" ]; then
    CLUSTER_IP=`kubectl get nodes -o json | jq -r '.items[0].spec.externalID'`
    CLUSTER_PORT=`kubectl get svc $1 -o json | jq -r '.spec.ports[0].nodePort'`
  fi
}

export_cluster_env () {
  export CLUSTER_TYPE
  export CLUSTER_IP
  export CLUSTER_PORT
  export CLUSTER_SERVER=http://$CLUSTER_IP:$CLUSTER_PORT

  env | grep "^CLUSTER_"
}

is_cluster_ready () {
  [ $(kubectl get ep $1 -o json | jq -r '.subsets | length') -gt 0 ]
}

wait_for_cluster () {
  case "$CLUSTER_TYPE" in
  "NodePort" )
    i=0
    while ! $(is_cluster_ready $1) && [ $i -lt 6 ]; do
      kubectl get ep $1
      echo [$i] Waiting for $1 to be ready...
      sleep 5;
      i=$[$i+1]
    done
    if $(is_cluster_ready $1); then
      echo "READY"
    else
      echo "FAILED"
    fi
    ;;
  esac
}

command=$1
service=$2

case "$command" in
env )
  set_cluster_env $service
  export_cluster_env
  ;;
wait )
  set_cluster_env $service
  wait_for_cluster $service
  ;;
esac

