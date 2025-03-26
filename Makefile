# MIT License
# Copyright (c) 2025 Toni Liesche
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

ifneq ("$(wildcard $(CURDIR)/build.properties)","")
	include $(CURDIR)/build.properties
endif

include $(CURDIR)/make/scylladb.mk
include $(CURDIR)/make/versioning.mk

CONTAINER_SCYLLA=test-scylla-1
TABLE_NAME_FAAS=function_as_a_service_executions
SCYLLA_HOST=test-scylla-1
SCYLLA_PORT=8000

PROJECTS=function-as-a-service
DEBIANVER=bookworm
GOLANGVER=1.24.1

build-docker-rc-%: set-version-rc set-commit
	docker buildx build \
		--pull \
		--platform linux/amd64,linux/arm64 \
		--build-arg CACHEBUST=$(shell date +%s) \
		--build-arg BUILDVER=$(run.build.version) \
		--build-arg COMMIT=$(run.commit) \
		--build-arg GOLANGVER=$(GOLANGVER) \
		--build-arg DEBIANVER=$(DEBIANVER) \
		-t tliesche/$*:$(run.build.version) \
		$(if $(PUSH),--push,--load) \
 		docker/$*

build-docker-patch-%: set-version-patch set-commit
	docker buildx build \
		--pull \
		--platform linux/amd64,linux/arm64 \
		--build-arg CACHEBUST=$(shell date +%s) \
		--build-arg BUILDVER=$(run.build.version) \
		--build-arg COMMIT=$(run.commit) \
		--build-arg GOLANGVER=$(GOLANGVER) \
		--build-arg DEBIANVER=$(DEBIANVER) \
		-t tliesche/$*:$(run.build.version) \
		-t tliesche/$*:$(run.build.version.minor) \
		-t tliesche/$*:$(run.build.version.major) \
		-t tliesche/$*:latest \
		$(if $(PUSH),--push,--load) \
		docker/$*

build-docker-release-%: set-version-release set-commit
	docker buildx build \
		--pull \
		--platform linux/amd64,linux/arm64 \
		--build-arg CACHEBUST=$(shell date +%s) \
		--build-arg BUILDVER=$(run.build.version) \
		--build-arg COMMIT=$(run.commit) \
		--build-arg GOLANGVER=$(GOLANGVER) \
		--build-arg DEBIANVER=$(DEBIANVER) \
		-t tliesche/$*:$(run.build.version) \
		-t tliesche/$*:$(run.build.version.minor) \
		-t tliesche/$*:$(run.build.version.major) \
		-t tliesche/$*:latest \
		$(if $(PUSH),--push,--load) \
		docker/$*

build-docker-%: set-version-%
	$(foreach project, $(PROJECTS), $(MAKE) build-docker-$*-$(project);)

build-dev-docker-faas:
	$(MAKE) build-dev-docker-function-as-a-service

build-dev-docker-%:
	docker build \
		--pull \
		--build-arg CACHEBUST=$(shell date +%s) \
		--build-arg BUILDVER=develop \
		--build-arg COMMIT=develop \
		--build-arg GOLANGVER=$(GOLANGVER) \
		--build-arg DEBIANVER=$(DEBIANVER) \
		-t tliesche/$*:develop \
		-f docker/$*/dev/Dockerfile \
		.

build-dev-docker:
	$(foreach project, $(PROJECTS_DEV), $(MAKE) build-dev-docker-$(project);)

build-%:
	go build -o build/$* cmd/$*/main.go

.PHONY: build
build:
	$(foreach project, $(PROJECTS), $(MAKE) build-$(project);)

down-test:
	docker compose -f docker/docker-compose.yml -p test down --volumes

up-test: up-docker create-scylla-tables

up-docker:
	docker compose -f docker/docker-compose.yml -p test up -d --remove-orphans
