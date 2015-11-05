SRCROOT ?= $(realpath .)

MACHINE_DEFAULT_IP = $(shell docker-machine ip default)

start: testenv-start

stop: testenv-stop

kubernetes-start:
	docker run --net=host -d gcr.io/google_containers/etcd:2.0.12 /usr/local/bin/etcd --addr=127.0.0.1:4001 --bind-addr=0.0.0.0:4001 --data-dir=/var/etcd/data
	docker run --volume=/:/rootfs:ro --volume=/sys:/sys:ro --volume=/dev:/dev --volume=/var/lib/docker/:/var/lib/docker:rw --volume=/var/lib/kubelet/:/var/lib/kubelet:rw --volume=/var/run:/var/run:rw --net=host --pid=host --privileged=true -d johnnylai/hyperkube:v1.0.6 /hyperkube kubelet --containerized=true --hostname-override="127.0.0.1" --address="0.0.0.0" --api-servers=http://localhost:8080 --config=/etc/kubernetes/manifests --v=5
	docker run -d --net=host --privileged johnnylai/hyperkube:v1.0.6 /hyperkube proxy --master=http://127.0.0.1:8080 --v=2

kubernetes-stop:
	-docker kill `docker ps | grep "gcr.io/google_containers" | awk '{print $$1}'`

testenv-start:
	kubectl create -f $(SRCROOT)/app-controller.yml
	kubectl create -f $(SRCROOT)/app-service.yml
	kubectl create -f $(SRCROOT)/db-controller.yml
	kubectl create -f $(SRCROOT)/db-service.yml

testenv-stop:
	kubectl delete all -lapp=go-service-basic

killall:
	-docker ps -a | sed '1d' | awk '{print $$1}' | xargs docker kill
  
forward-ports:
	ssh -i  ~/.docker/machine/machines/default/id_rsa docker@$(MACHINE_DEFAULT_IP) -L8080:localhost:8080
