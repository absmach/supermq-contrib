# Copyright (c) Abstract Machines
# SPDX-License-Identifier: Apache-2.0

SMQ_DOCKER_IMAGE_NAME_PREFIX ?= supermq-contrib
BUILD_DIR = build
SERVICES = opcua lora influxdb-writer influxdb-reader mongodb-writer \
	mongodb-reader cassandra-writer cassandra-reader smtp-notifier smpp-notifier
TEST_API_SERVICES = notifiers readers twins
TEST_API = $(addprefix test_api_,$(TEST_API_SERVICES))
DOCKERS = $(addprefix docker_,$(SERVICES))
DOCKERS_DEV = $(addprefix docker_dev_,$(SERVICES))
CGO_ENABLED ?= 0
GOARCH ?= amd64
VERSION ?= $(shell git describe --abbrev=0 --tags 2>/dev/null || echo 'unknown')
COMMIT ?= $(shell git rev-parse HEAD)
TIME ?= $(shell date +%F_%T)
USER_REPO ?= $(shell git remote get-url origin | sed -e 's/.*\/\([^/]*\)\/\([^/]*\).*/\1_\2/' )
empty:=
space:= $(empty) $(empty)
# Docker compose project name should follow this guidelines: https://docs.docker.com/compose/reference/#use--p-to-specify-a-project-name
DOCKER_PROJECT ?= $(shell echo $(subst $(space),,$(USER_REPO)) | tr -c -s '[:alnum:][=-=]' '_' | tr '[:upper:]' '[:lower:]')
DOCKER_COMPOSE_COMMANDS_SUPPORTED := up down config
DEFAULT_DOCKER_COMPOSE_COMMAND  := up
GRPC_MTLS_CERT_FILES_EXISTS = 0
MOCKERY_VERSION=v3.5.0
ifneq ($(SMQ_MESSAGE_BROKER_TYPE),)
	SMQ_MESSAGE_BROKER_TYPE := $(SMQ_MESSAGE_BROKER_TYPE)
else
	SMQ_MESSAGE_BROKER_TYPE=msg_nats
endif

ifneq ($(SMQ_ES_TYPE),)
	SMQ_ES_TYPE := $(SMQ_ES_TYPE)
else
	SMQ_ES_TYPE=es_nats
endif

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) \
	go build -tags $(SMQ_MESSAGE_BROKER_TYPE) --tags $(SMQ_ES_TYPE) -ldflags "-s -w \
	-X 'github.com/absmach/supermq.BuildTime=$(TIME)' \
	-X 'github.com/absmach/supermq.Version=$(VERSION)' \
	-X 'github.com/absmach/supermq.Commit=$(COMMIT)'" \
	-o ${BUILD_DIR}/$(1) cmd/$(1)/main.go
endef

define make_docker
	$(eval svc=$(subst docker_,,$(1)))

	docker build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--build-arg GOARCH=$(GOARCH) \
		--build-arg GOARM=$(GOARM) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg TIME=$(TIME) \
		--tag=$(SMQ_DOCKER_IMAGE_NAME_PREFIX)/$(svc) \
		-f docker/Dockerfile .
endef

define make_docker_dev
	$(eval svc=$(subst docker_dev_,,$(1)))

	docker build \
		--no-cache \
		--build-arg SVC=$(svc) \
		--tag=$(SMQ_DOCKER_IMAGE_NAME_PREFIX)/$(svc) \
		-f docker/Dockerfile.dev ./build
endef

ifneq ("$(wildcard docker/ssl/certs/*-grpc-*)","")
GRPC_MTLS_CERT_FILES_EXISTS = 1
else
GRPC_MTLS_CERT_FILES_EXISTS = 0
endif

all: $(SERVICES)

.PHONY: all $(SERVICES) dockers dockers_dev latest release run run_addons grpc_mtls_certs check_mtls check_certs test_api

clean:
	rm -rf ${BUILD_DIR}

cleandocker:
	# Stops containers and removes containers, networks, volumes, and images created by up
	docker compose -f docker/docker-compose.yml -p $(DOCKER_PROJECT) down --rmi all -v --remove-orphans

ifdef pv
	# Remove unused volumes
	docker volume ls -f name=$(SMQ_DOCKER_IMAGE_NAME_PREFIX) -f dangling=true -q | xargs -r docker volume rm
endif

install:
	for file in $(BUILD_DIR)/*; do \
		cp $$file $(GOBIN)/supermq-contrib-`basename $$file`; \
	done

mocks:
	@which mockery > /dev/null || go install github.com/vektra/mockery/v3@$(MOCKERY_VERSION)
	@unset MOCKERY_VERSION 
	mockery --config ./tools/config/mockery.yaml

DIRS = consumers readers opcua twins lora
test: mocks
	mkdir -p coverage
	@for dir in $(DIRS); do \
        go test -v --race -count 1 -tags test -coverprofile=coverage/$$dir.out $$(go list ./... | grep $$dir | grep -v 'cmd'); \
    done

define test_api_service
	$(eval svc=$(subst test_api_,,$(1)))
	@which st > /dev/null || (echo "schemathesis not found, please install it from https://github.com/schemathesis/schemathesis#getting-started" && exit 1)

	@if [ -z "$(USER_TOKEN)" ]; then \
		echo "USER_TOKEN is not set"; \
		echo "Please set it to a valid token"; \
		exit 1; \
	fi

	@if [ "$(svc)" = "http" ] && [ -z "$(CLIENT_SECRET)" ]; then \
		echo "CLIENT_SECRET is not set"; \
		echo "Please set it to a valid secret"; \
		exit 1; \
	fi

	st run api/openapi/$(svc).yml \
	--checks all \
	--base-url $(2) \
	--header "Authorization: Bearer $(USER_TOKEN)" \
	--contrib-openapi-formats-uuid \
	--hypothesis-suppress-health-check=filter_too_much \
	--stateful=links;
endef

test_api_twins: TEST_API_URL := http://localhost:9018
test_api_readers: TEST_API_URL := http://localhost:9009 # This can be the URL of any reader service.
test_api_notifiers: TEST_API_URL := http://localhost:9014 # This can be the URL of any notifier service.

$(TEST_API):
	$(call test_api_service,$(@),$(TEST_API_URL))

$(SERVICES):
	$(call compile_service,$(@))

$(DOCKERS):
	$(call make_docker,$(@),$(GOARCH))

$(DOCKERS_DEV):
	$(call make_docker_dev,$(@))

dockers: $(DOCKERS)
dockers_dev: $(DOCKERS_DEV)

define docker_push
	for svc in $(SERVICES); do \
		docker push $(SMQ_DOCKER_IMAGE_NAME_PREFIX)/$$svc:$(1); \
	done
endef

changelog:
	git log $(shell git describe --tags --abbrev=0)..HEAD --pretty=format:"- %s"

latest: dockers
	$(call docker_push,latest)

release:
	$(eval version = $(shell git describe --abbrev=0 --tags))
	git checkout $(version)
	$(MAKE) dockers
	for svc in $(SERVICES); do \
		docker tag $(SMQ_DOCKER_IMAGE_NAME_PREFIX)/$$svc $(SMQ_DOCKER_IMAGE_NAME_PREFIX)/$$svc:$(version); \
	done
	$(call docker_push,$(version))

rundev:
	cd scripts && ./run.sh
