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

create-scylla-tables:
	@echo "Checking if the Scylla container is running..."
	@RETRY_LIMIT=10; \
	RETRY_COUNT=0; \
	until docker inspect -f '{{.State.Running}}' $(SCYLLA_CONTAINER) 2>/dev/null | grep -q "true" || [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; do \
		echo "Scylla container is not running yet. Retrying..."; \
		RETRY_COUNT=$$((RETRY_COUNT + 1)); \
		sleep 2; \
	done; \
	if [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; then \
		echo "Timeout reached, Scylla container is not running."; \
		exit 1; \
	fi
	@echo "Scylla container is running. Proceeding with Alternator API setup..."

	@echo "Checking if Scylla Alternator API is available inside the container..."
	@RETRY_LIMIT=10; \
	RETRY_COUNT=0; \
	docker exec -it $(SCYLLA_CONTAINER) /bin/bash -c "\
		RETRY_LIMIT=10; \
		RETRY_COUNT=0; \
		until curl -s --get http://$(SCYLLA_HOST):$(SCYLLA_PORT) | grep -q 'healthy: $(SCYLLA_HOST):$(SCYLLA_PORT)' || [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; do \
			echo 'Scylla Alternator API not available yet. Retrying...'; \
			RETRY_COUNT=$$((RETRY_COUNT + 1)); \
			sleep 2; \
		done; \
		if [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; then \
			echo 'Timeout reached, Scylla Alternator API is not running.'; \
			exit 1; \
		fi; \
		echo 'Scylla Alternator API is up!'; \
		echo 'Check if table \"$(TABLE_NAME_FAAS)\" already exists'; \
		if curl -s http://$(SCYLLA_HOST):$(SCYLLA_PORT) \
			-H 'Content-Type: application/json' \
			-H 'X-Amz-Target: DynamoDB_20120810.DescribeTable' \
			-d '{\"TableName\": \"$(TABLE_NAME_FAAS)\"}' | grep -q 'ResourceNotFoundException'; then \
			echo 'Table does not exist, creating table \"$(TABLE_NAME_FAAS)\"...'; \
			curl -o /dev/null -s -f -X POST http://$(SCYLLA_HOST):$(SCYLLA_PORT) \
			    -H 'Content-Type: application/json' \
				-H 'X-Amz-Target: DynamoDB_20120810.CreateTable' \
			    -d '{\
			          \"TableName\": \"$(TABLE_NAME_FAAS)\", \
			          \"KeySchema\": [{\"AttributeName\": \"id\", \"KeyType\": \"HASH\"}], \
			          \"AttributeDefinitions\": [{\"AttributeName\": \"id\", \"AttributeType\": \"S\"}], \
			          \"ProvisionedThroughput\": {\"ReadCapacityUnits\": 1, \"WriteCapacityUnits\": 1} \
			     }'   && echo 'Success: Table creation request sent!' && \
			curl -o /dev/null -s -f -X POST http://$(SCYLLA_HOST):$(SCYLLA_PORT) \
            	-H 'Content-Type: application/json' \
            	-H 'X-Amz-Target: DynamoDB_20120810.UpdateTimeToLive' \
            	-d '{\
                      \"TableName\": \"$(TABLE_NAME_FAAS)\", \
    				  \"TimeToLiveSpecification\": { \
        			  \"Enabled\": true, \
        			  \"AttributeName\": \"ttl\" \
    				} \
                 }'   && echo 'Success: TTL enabled!' \
            || echo 'Error: Table creation request failed!'; \
		else \
			echo 'Table \"$(TABLE_NAME_FAAS)\" already exists. Skipping table creation.'; \
		fi"
