NETWORKNAME = test-rabbit-network
RABBIT = rabbit
MONGO = mongo
CONSUL = consul

NAME = worker
DEFAULT_PATH = /
CONSUL_ADDR = $(CONSUL):8500

MONITOR_NAME = monitor
HTTP_ADDR_MONITOR = :8082
HTTP_ADDR_MONITOR_P = 8082

start_worker:
	docker run --network=$(NETWORKNAME) -it --rm --hostname=$(NAME) --name $(NAME) \
		-e ROOT=$(DEFAULT_PATH) -e SERVER_NAME=$(NAME) -e NETWORKNAME=$(RABBIT) \
		-e CONSUL_ADDR=$(CONSUL_ADDR) opbx_worker	

build_worker:
	docker build -t opbx_worker ./worker

start_monitor:
	docker run --network=$(NETWORKNAME) -it --rm -p $(HTTP_ADDR_MONITOR_P):$(HTTP_ADDR_MONITOR_P) \
		--hostname=$(MONITOR_NAME) --name $(MONITOR_NAME) -e NETWORKNAME=$(RABBIT) \
		-e HTTP_ADDR=$(HTTP_ADDR_MONITOR) opbx_monitor

build_monitor:
	docker build -t opbx_monitor ./monitor

start_services:
	docker run --name=$(RABBIT) -host=$(RABBIT) --network $(NETWORKNAME) --rm -d rabbitmq:latest
	docker run --name=$(MONGO) -host=$(MONGO) --network $(NETWORKNAME) --rm -d mongo:latest
	docker run --network=$(NETWORKNAME)  --name=$(CONSUL) --hostname=$(CONSUL) -p 8500:8500 \
		--rm -d -e CONSUL_LOCAL_CONFIL='{"server":true, "enable_debug":true}' \
		consul:latest agent -dev -ui -client=0.0.0.0
	docker ps

pull_services:
	docker pull rabbitmq:latest
	docker pull mongo:latest
	docker pull consul:latest

stop_services:
	docker stop $(RABBIT) && docker stop $(MONGO) && docker stop $(CONSUL) 
