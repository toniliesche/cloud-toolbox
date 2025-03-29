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

QUIET= --quiet

setup-rabbitmq:
	$(MAKE) create-rabbitmq-default-users

create-rabbitmq-default-users:
	$(call print_message,Creating RabbitMQ User "$(RABBITMQ_ADMIN_USER)")
	$(MAKE) create-rabbitmq-user RABBITMQ_USER=$(RABBITMQ_ADMIN_USER) RABBITMQ_PASSWORD=$(RABBITMQ_ADMIN_PASSWORD)

	$(call print_message,Setting RabbitMQ User "$(RABBITMQ_ADMIN_USER)" Tags)
	$(MAKE) set-rabbitmq-user-tags RABBITMQ_USER=$(RABBITMQ_ADMIN_USER) RABBITMQ_TAG=administrator

	$(call print_message,Setting RabbitMQ Permissions for User "$(RABBITMQ_ADMIN_USER)" on vHost "/")
	$(MAKE) set-rabbitmq-user-permissions RABBITMQ_USER=$(RABBITMQ_ADMIN_USER) RABBITMQ_VHOST=/ RABBITMQ_PERM_CONF=".*" RABBITMQ_PERM_WRITE=".*" RABBITMQ_PERM_READ=".*"

	docker exec -u 100 $(RABBITMQ_CONTAINER) /bin/bash -c "\
		echo '[default]' > ~/.rabbitmqadmin.conf \
		&& echo 'hostname=localhost' >> ~/.rabbitmqadmin.conf \
		&& echo 'port=15672' >> ~/.rabbitmqadmin.conf \
		&& echo 'username=$(RABBITMQ_ADMIN_USER)' >> ~/.rabbitmqadmin.conf \
		&& echo 'password=$(RABBITMQ_ADMIN_PASSWORD)' >> ~/.rabbitmqadmin.conf"

	$(call print_message,Deleting RabbitMQ User "guest")
	$(MAKE) delete-rabbitmq-user RABBITMQ_USER=guest

check-rabbitmq-container:
	echo "Checking if the RabbitMQ container is running..."
	RETRY_LIMIT=10; \
	RETRY_COUNT=0; \
	until docker inspect -f '{{.State.Running}}' $(RABBITMQ_CONTAINER) 2>/dev/null | grep -q "true" || [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; do \
		echo "RabbitMQ container is not running yet. Retrying... ($$RETRY_COUNT/$$RETRY_LIMIT)"; \
		RETRY_COUNT=$$((RETRY_COUNT + 1)); \
		sleep 2; \
	done; \
	if [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; then \
		echo "Timeout reached, RabbitMQ container is not running."; \
		exit 1; \
	fi
	echo "RabbitMQ container is running."

check-rabbitmq: check-rabbitmq-container
	@echo "Checking if RabbitMQ is ready..."
	@RETRY_LIMIT=10; \
	RETRY_COUNT=0; \
	until docker exec -u 100 $(RABBITMQ_CONTAINER) rabbitmqctl status >/dev/null 2>&1 || [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; do \
		echo "RabbitMQ is not ready yet. Retrying... ($$RETRY_COUNT/$$RETRY_LIMIT)"; \
		RETRY_COUNT=$$((RETRY_COUNT + 1)); \
		sleep 2; \
	done; \
	if [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; then \
		echo "Timeout reached. RabbitMQ is not running or not responding."; \
		exit 1; \
	fi; \
	echo "RabbitMQ is up and ready!"

check-rabbitmq-vhost: check-rabbitmq
	echo "Checking if RabbitMQ vHost "$(RABBITMQ_VHOST)" exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "rabbitmqctl list_vhosts 2> /dev/null | grep -q -w '$(RABBITMQ_VHOST)'; \
		if [ $$? -eq 0 ]; then \
			echo 'RabbitMQ vHost \"$(RABBITMQ_VHOST)\" does exist.'; \
		else \
			echo 'Error: RabbitMQ vHost \"$(RABBITMQ_VHOST)\" does not exist!'; \
			exit 1; \
		fi"

create-rabbitmq-vhost: check-rabbitmq
	echo "Creating RabbitMQ vHost \"$(RABBITMQ_VHOST)\"..."
	echo "Checking if RabbitMQ vHost \"$(RABBITMQ_VHOST)\" already exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_vhosts 2> /dev/null | grep -q -w '$(RABBITMQ_VHOST)'; then \
			echo 'RabbitMQ vHost \"$(RABBITMQ_VHOST)\" does already exist, skipping vHost creation.'; \
		else \
			echo 'RabbitMQ vHost \"$(RABBITMQ_VHOST)\" does not exist, running vHost creation now.'; \
			rabbitmqctl $(QUIET) add_vhost $(RABBITMQ_VHOST) \
			&& echo 'Success: Finished creating RabbitMQ vHost \"$(RABBITMQ_VHOST)\".' \
			|| echo 'Error: Failed creating RabbitMQ vHost \"$(RABBITMQ_VHOST)\"!'; \
		fi"

delete-rabbitmq-vhost: check-rabbitmq
	echo "Deleting RabbitMQ vHost \"$(RABBITMQ_VHOST)\"..."
	echo "Checking if RabbitMQ vHost \"$(RABBITMQ_VHOST)\" exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_vhosts 2> /dev/null | grep -q -w '$(RABBITMQ_VHOST)'; then \
			echo 'RabbitMQ vHost \"$(RABBITMQ_VHOST)\" does exist, running vHost deletion now.'; \
			rabbitmqctl $(QUIET) delete_vhost $(RABBITMQ_VHOST) \
			&& echo 'Success: Finished deleting RabbitMQ vHost \"$(RABBITMQ_VHOST)\".' \
			|| echo 'Error: Failed deleting RabbitMQ vHost \"$(RABBITMQ_VHOST)\"!'; \
		else \
			echo 'RabbitMQ vHost \"$(RABBITMQ_VHOST)\" does not exist, skipping vHost deletion.'; \
		fi"

check-rabbitmq-exchange: check-rabbitmq check-rabbitmq-vhost
	echo "Checking if RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\" exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_exchanges -p $(RABBITMQ_VHOST) 2> /dev/null | grep -q -E '$(RABBITMQ_EXCHANGE)[[:space:]]+'; then \
			echo 'RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\" does exist.'; \
		else \
			echo 'Error: RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\" does not exist!'; \
			exit 1; \
		fi"

create-rabbitmq-exchange: check-rabbitmq check-rabbitmq-vhost
	echo "Creating RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\"..."
	echo "Checking if RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\" already exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_exchanges -p $(RABBITMQ_VHOST) 2> /dev/null | grep -q -E '$(RABBITMQ_EXCHANGE)[[:space:]]+'; then \
			echo 'RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\" does already exist, skipping exchange creation.'; \
		else \
			echo 'RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\" does not exist, running exchange creation now.'; \
			rabbitmqadmin $(QUIET) declare exchange name=$(RABBITMQ_EXCHANGE) type=direct durable=true --vhost=$(RABBITMQ_VHOST) \
			&& echo 'Success: Finished creating RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\".' \
			|| echo 'Error: Failed creating RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\"!'; \
		fi"

delete-rabbitmq-exchange: check-rabbitmq check-rabbitmq-vhost
	echo "Deleting RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\"..."
	echo "Checking if RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\" exists..."
	docker exec $(RABBITMQ_EXCHANGE) /bin/bash -c "\
		if rabbitmqctl list_exchanges -p $(RABBITMQ_VHOST) 2> /dev/null | grep -q -E '$(RABBITMQ_EXCHANGE)[[:space:]]+'; then \
			echo 'RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\" does exist, running exchange deletion now.'; \
			rabbitmqadmin $(QUIET) delete exchange name=$(RABBITMQ_EXCHANGE) --vhost=$(RABBITMQ_VHOST) \
			&& echo 'Success: Finished deleting RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\".' \
			|| echo 'Error: Failed deleting RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\"!'; \
		else \
			echo 'RabbitMQ Exchange \"$(RABBITMQ_EXCHANGE)\" on vHost \"$(RABBITMQ_VHOST)\" does not exist, skipping exchange deletion.'; \
		fi"

check-rabbitmq-queue: check-rabbitmq check-rabbitmq-vhost
	echo "Checking if RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\" exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_queues -p $(RABBITMQ_VHOST) 2> /dev/null | grep -q -E '$(RABBITMQ_QUEUE)[[:space:]]+'; then \
			echo 'RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\" does exist.'; \
		else \
			echo 'Error: RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\" does not exist!'; \
			exit 1; \
		fi"

create-rabbitmq-queue: check-rabbitmq check-rabbitmq-vhost
	echo "Creating RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\"..."
	echo "Checking if RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\" already exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_queues -p $(RABBITMQ_VHOST) 2> /dev/null | grep -q -E '$(RABBITMQ_QUEUE)[[:space:]]+'; then \
			echo 'RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\" does already exist, skipping queue creation.'; \
		else \
			echo 'RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\" does not exist, running queue creation now.'; \
			QUEUE_ARGS='durable=true'; \
			if [ -n \"$(RABBITMQ_DLX)\" ]; then \
				echo 'Configuring Dead Letter Exchange to \"$(RABBITMQ_DLX)\"'; \
				QUEUE_ARGS=\"\$$QUEUE_ARGS arguments={\\\"x-dead-letter-exchange\\\":\\\"$(RABBITMQ_DLX)\\\",\\\"x-dead-letter-routing-key\\\":\\\"$(RABBITMQ_QUEUE).dlq\\\",\\\"x-max-delivery-count\\\":3}\"; \
			fi; \
			rabbitmqadmin $(QUIET) declare queue name=$(RABBITMQ_QUEUE) \$$QUEUE_ARGS --vhost=$(RABBITMQ_VHOST) \
				&& echo 'Success: Finished creating RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\".' \
				|| echo 'Error: Failed creating RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\"!'; \
		fi"

delete-rabbitmq-queue: check-rabbitmq check-rabbitmq-vhost
	echo "Deleting RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\"..."
	echo "Checking if RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\" exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_queues -p $(RABBITMQ_VHOST) 2> /dev/null | grep -q -E '$(RABBITMQ_QUEUE)[[:space:]]+'; then \
			echo 'RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\" does exist, running queue deletion now.'; \
			rabbitmqadmin $(QUIET) delete queue name=$(RABBITMQ_QUEUE) --vhost=$(RABBITMQ_VHOST) \
			&& echo 'Success: Finished deleting RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\".' \
			|| echo 'Error: Failed deleting RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\"!'; \
		else \
			echo 'RabbitMQ Queue \"$(RABBITMQ_QUEUE)\" on vHost \"$(RABBITMQ_VHOST)\" does not exist, skipping queue deletion.'; \
		fi"

create-rabbitmq-binding: check-rabbitmq check-rabbitmq-vhost check-rabbitmq-queue check-rabbitmq-exchange
	echo "Creating RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\")..."
	echo "Checking if RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\") already exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_bindings -p $(RABBITMQ_VHOST) 2> /dev/null | grep -q -E '$(RABBITMQ_EXCHANGE)[[:space:]]+exchange[[:space:]]+$(RABBITMQ_QUEUE)[[:space:]]+queue[[:space:]]+$(RABBITMQ_BINDING_KEY)[[:space:]]+'; then \
			echo 'RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\") does already exist, skipping binding creation.'; \
		else \
			echo 'RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\") does not exist, running binding creation now.'; \
			rabbitmqadmin $(QUIET) declare binding source=$(RABBITMQ_EXCHANGE) destination=$(RABBITMQ_QUEUE) routing_key=$(RABBITMQ_BINDING_KEY) --vhost=$(RABBITMQ_VHOST) \
			&& echo 'Success: Finished creating RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\").' \
			|| echo 'Error: Failed creating RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\")!'; \
		fi"

delete-rabbitmq-binding: check-rabbitmq check-rabbitmq-vhost
	echo "Deleting RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\")..."
	echo "Checking if RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\") exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_bindings -p $(RABBITMQ_VHOST) 2> /dev/null | grep -q -E '$(RABBITMQ_EXCHANGE)[[:space:]]+exchange[[:space:]]+$(RABBITMQ_QUEUE)[[:space:]]+queue[[:space:]]+$(RABBITMQ_BINDING_KEY)[[:space:]]+'; then \
			echo 'RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\") does exist, running binding deletion now.'; \
			rabbitmqadmin $(QUIET) delete binding source=$(RABBITMQ_EXCHANGE) destination=$(RABBITMQ_QUEUE) routing_key=$(RABBITMQ_BINDING_KEY) --vhost=$(RABBITMQ_VHOST) \
			&& echo 'Success: Finished deleting RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\").' \
			|| echo 'Error: Failed deleting RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\")!'; \
		else \
			echo 'RabbitMQ Binding from Exchange \"$(RABBITMQ_EXCHANGE)\" to Queue \"$(RABBITMQ_QUEUE)\" (Topic: \"$(RABBITMQ_BINDING_KEY)\") does not exist, skipping binding deletion.'; \
		fi"

check-rabbitmq-user: check-rabbitmq
	echo "Checking if RabbitMQ User "$(RABBITMQ_USER)" exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "rabbitmqctl list_users 2> /dev/null | grep -q -w '$(RABBITMQ_USER)'; \
		if [ $$? -eq 0 ]; then \
			echo 'RabbitMQ User \"$(RABBITMQ_USER)\" does exist.'; \
		else \
			echo 'Error: RabbitMQ User \"$(RABBITMQ_USER)\" does not exist!'; \
			exit 1; \
		fi"

create-rabbitmq-user: check-rabbitmq
	echo "Creating RabbitMQ User \"$(RABBITMQ_USER)\"..."
	echo "Checking if RabbitMQ User \"$(RABBITMQ_USER)\" already exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_users 2> /dev/null | grep -q -w '$(RABBITMQ_USER)'; then \
			echo 'RabbitMQ User \"$(RABBITMQ_USER)\" does already exist, skipping user creation.'; \
		else \
			echo 'RabbitMQ User \"$(RABBITMQ_USER)\" does not exist, running user creation now.'; \
			rabbitmqctl $(QUIET) add_user $(RABBITMQ_USER) $(RABBITMQ_PASSWORD) \
			&& echo 'Success: Finished creating RabbitMQ User \"$(RABBITMQ_USER)\".' \
			|| echo 'Error: Failed creating RabbitMQ User \"$(RABBITMQ_USER)\"!'; \
		fi"

delete-rabbitmq-user: check-rabbitmq
	echo "Deleting RabbitMQ User \"$(RABBITMQ_USER)\"..."
	echo "Checking if RabbitMQ User \"$(RABBITMQ_USER)\" exists..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		if rabbitmqctl list_users 2> /dev/null | grep -q -w '$(RABBITMQ_USER)'; then \
			echo 'RabbitMQ User \"$(RABBITMQ_USER)\" does exist, running user deletion now.'; \
			rabbitmqctl $(QUIET) delete_user $(RABBITMQ_USER) \
			&& echo 'Success: Finished deleting RabbitMQ User \"$(RABBITMQ_USER)\".' \
			|| echo 'Error: Failed deleting RabbitMQ User \"$(RABBITMQ_USER)\"!'; \
		else \
			echo 'RabbitMQ User \"$(RABBITMQ_USER)\" does not exist, skipping user deletion.'; \
		fi"

set-rabbitmq-user-permissions: check-rabbitmq check-rabbitmq-user check-rabbitmq-vhost
	echo "Setting Permissions for RabbitMQ User \"$(RABBITMQ_USER)\" on vHost \"$(RABBITMQ_VHOST)\"..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		rabbitmqctl $(QUIET) set_permissions -p $(RABBITMQ_VHOST) $(RABBITMQ_USER) '$(RABBITMQ_PERM_CONF)' '$(RABBITMQ_PERM_WRITE)' '$(RABBITMQ_PERM_READ)' \
		&& echo 'Success: Finished setting Permissions for RabbitMQ User \"$(RABBITMQ_USER)\".' \
		|| echo 'Error: Failed setting Permissions for RabbitMQ User \"$(RABBITMQ_USER)\"!';"

set-rabbitmq-user-tags: check-rabbitmq check-rabbitmq-user
	echo "Setting Tag \"$(RABBITMQ_TAG)\" for RabbitMQ User \"$(RABBITMQ_USER)\"..."
	docker exec $(RABBITMQ_CONTAINER) /bin/bash -c "\
		rabbitmqctl $(QUIET) set_user_tags $(RABBITMQ_USER) $(RABBITMQ_TAG) \
		&& echo 'Success: Finished setting Tag \"$(RABBITMQ_TAG)\" for RabbitMQ User \"$(RABBITMQ_USER)\".' \
		|| echo 'Error: Failed setting Tag \"$(RABBITMQ_TAG)\" for RabbitMQ User \"$(RABBITMQ_USER)\"!';"
