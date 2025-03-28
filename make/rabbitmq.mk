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

create-rabbitmq-admin:
	@echo "Checking if the RabbitMQ container is running..."
	@RETRY_LIMIT=10; \
	RETRY_COUNT=0; \
	until docker inspect -f '{{.State.Running}}' $(RABBITMQ_CONTAINER) 2>/dev/null | grep -q "true" || [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; do \
		echo "RabbitMQ container is not running yet. Retrying..."; \
		RETRY_COUNT=$$((RETRY_COUNT + 1)); \
		sleep 2; \
	done; \
	if [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; then \
		echo "Timeout reached, RabbitMQ container is not running."; \
		exit 1; \
	fi
	@echo "RabbitMQ container is running. Proceeding with RabbitMQ admin user setup..."

	@echo "Checking if RabbitMQ admin user is available inside the container..."
	@RETRY_LIMIT=10; \
	RETRY_COUNT=0; \
	docker exec -it $(RABBITMQ_CONTAINER) /bin/bash -c "\
		RETRY_LIMIT=10; \
		RETRY_COUNT=0; \
		until rabbitmqctl list_users | grep -q '$(RABBITMQ_USER)' || [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; do \
			echo 'RabbitMQ admin user not available yet. Retrying...'; \
			RETRY_COUNT=$$((RETRY_COUNT + 1)); \
			sleep 2; \
		done; \
		if [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; then \
			echo 'Timeout reached, RabbitMQ admin user is not available.'; \
			exit 1; \
		fi; \
		echo 'RabbitMQ admin user is up!'; \
		echo 'Check if user \"$(RABBITMQ_USER)\" already exists'; \
		if rabbitmqctl list_users | grep -q -w '$(RABBITMQ_USER)'; then \
			echo 'User already exists, skipping user creation'; \
		else \
			echo 'User does not exist, creating user \"$(RABBITMQ_USER)\"...'; \
			rabbitmqctl add_user $(RABBITMQ_USER) $(RABBITMQ_PASSWORD) && echo 'Success: User creation request sent!'; \
			rabbitmqctl set_user_tags $(RABBITMQ_USER) administrator && echo 'Success: User tags set!'; \
			rabbitmqctl set_permissions -p / $(RABBITMQ_USER) '.*' '.*' '.*' && echo 'Success: User permissions set'; \
		fi"