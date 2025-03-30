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

MAKEFLAGS += --no-print-directory -s

ifneq ("$(wildcard $(CURDIR)/build.properties)","")
	include $(CURDIR)/build.properties
endif

ifneq ("$(wildcard $(CURDIR)/.env)","")
	include $(CURDIR)/.env
endif

include $(CURDIR)/make/functions.mk
include $(CURDIR)/make/faas.mk
include $(CURDIR)/make/ft.mk
include $(CURDIR)/make/rabbitmq.mk
include $(CURDIR)/make/scylladb.mk
include $(CURDIR)/make/versioning.mk

build-docker-rc-%: set-version-rc set-commit
	$(call print_message,Building release candidate Docker image for "$*")

	echo "Build version: $(run.build.version)"
	echo "Build commit: $(run.commit)"
	echo "Golang version: $(GOLANGVER)"
	echo "Debian version: $(DEBIANVER)"
	echo

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
	$(call print_message,Building patch Docker image for "$*")

	echo "Build version: $(run.build.version)"
	echo "Build commit: $(run.commit)"
	echo "Golang version: $(GOLANGVER)"
	echo "Debian version: $(DEBIANVER)"
	echo

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
	$(call print_message,Building release Docker image for "$*")

	echo "Build version: $(run.build.version)"
	echo "Build commit: $(run.commit)"
	echo "Golang version: $(GOLANGVER)"
	echo "Debian version: $(DEBIANVER)"
	echo

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

build-dev-docker-ft:
	$(MAKE) build-dev-docker-function-trigger

build-dev-docker-%:
	$(call print_message,Building development Docker image for "$*")

	echo "Build version: develop"
	echo "Build commit: develop"
	echo "Golang version: $(GOLANGVER)"
	echo "Debian version: $(DEBIANVER)"
	echo

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
	$(call print_message,Stopping Docker containers)

	docker compose --env-file .env -f docker/docker-compose.yml -p test down --volumes

up-test: up-docker

setup-test: setup-rabbitmq setup-faas setup-ft

up-docker:
	$(call print_message,Starting Docker containers)

	docker compose --env-file .env -f docker/docker-compose.yml -p test up -d --remove-orphans
