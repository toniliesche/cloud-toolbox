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

check-scylla-container:
	echo "Checking if the Scylla container is running..."
	RETRY_LIMIT=10; \
	RETRY_COUNT=0; \
	until docker inspect -f '{{.State.Running}}' $(SCYLLA_CONTAINER) 2>/dev/null | grep -q "true" || [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; do \
		echo "Scylla container is not running yet. Retrying... ($$RETRY_COUNT/$$RETRY_LIMIT)"; \
		RETRY_COUNT=$$((RETRY_COUNT + 1)); \
		sleep 2; \
	done; \
	if [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; then \
		echo "Timeout reached, Scylla container is not running."; \
		exit 1; \
	fi
	echo "Scylla container is running."

check-scylla-api: check-scylla-container
	echo "Checking if Scylla Alternator API is available..."
	RETRY_LIMIT=10; \
	RETRY_COUNT=0; \
	docker exec $(SCYLLA_CONTAINER) /bin/bash -c "RETRY_LIMIT=10; \
		RETRY_COUNT=0; \
		until curl -s --get http://$(SCYLLA_HOST):$(SCYLLA_PORT) | grep -q 'healthy: $(SCYLLA_HOST):$(SCYLLA_PORT)' || [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; do \
			echo \"Scylla Alternator API not available yet. Retrying... ($$RETRY_COUNT/$$RETRY_LIMIT)\"; \
			RETRY_COUNT=$$((RETRY_COUNT + 1)); \
			sleep 2; \
		done; \
		if [ $$RETRY_COUNT -ge $$RETRY_LIMIT ]; then \
			echo 'Timeout reached, Scylla Alternator API is not running.'; \
			exit 1; \
		fi; \
		echo 'Scylla Alternator API is up and running.';"

create-scylla-table: check-scylla-api
	echo 'Check if table "$(SCYLLA_TABLE)" already exists';
	docker exec $(SCYLLA_CONTAINER) /bin/bash -c "if curl -s http://$(SCYLLA_HOST):$(SCYLLA_PORT) \
			-H 'Content-Type: application/json' \
			-H 'X-Amz-Target: DynamoDB_20120810.DescribeTable' \
			-d '{\"TableName\": \"$(SCYLLA_TABLE)\"}' | grep -q 'ResourceNotFoundException'; then \
			echo 'Table does not exist, creating table \"$(SCYLLA_TABLE)\"...'; \
			curl -o /dev/null -s -f -X POST http://$(SCYLLA_HOST):$(SCYLLA_PORT) \
			    -H 'Content-Type: application/json' \
				-H 'X-Amz-Target: DynamoDB_20120810.CreateTable' \
			    -d '{\
			          \"TableName\": \"$(SCYLLA_TABLE)\", \
			          \"KeySchema\": [{\"AttributeName\": \"id\", \"KeyType\": \"HASH\"}], \
			          \"AttributeDefinitions\": [{\"AttributeName\": \"id\", \"AttributeType\": \"S\"}], \
			          \"ProvisionedThroughput\": {\"ReadCapacityUnits\": 1, \"WriteCapacityUnits\": 1} \
			        }' \
			&& echo 'Success: Table creation request for table \"$(SCYLLA_TABLE)\" sent!' \
            || echo 'Error: Table creation request for table \"$(SCYLLA_TABLE)\" failed!'; \
		else \
			echo 'Table \"$(SCYLLA_TABLE)\" already exists. Skipping table creation.'; \
		fi"

describe-ttl:
	docker exec $(SCYLLA_CONTAINER) /bin/bash -c "curl -X POST http://$(SCYLLA_HOST):$(SCYLLA_PORT) \
            	-H 'Content-Type: application/json' \
            	-H 'X-Amz-Target: DynamoDB_20120810.DescribeTimeToLive' \
            	-d '{\
                      \"TableName\": \"$(SCYLLA_TABLE)\" \
                    }' && echo"

enable-ttl: check-scylla-api
	docker exec $(SCYLLA_CONTAINER) /bin/bash -c "curl -o /dev/null -s -f -X POST http://$(SCYLLA_HOST):$(SCYLLA_PORT) \
            	-H 'Content-Type: application/json' \
            	-H 'X-Amz-Target: DynamoDB_20120810.UpdateTimeToLive' \
            	-d '{\
                      \"TableName\": \"$(SCYLLA_TABLE)\", \
    				  \"TimeToLiveSpecification\": { \
        			    \"Enabled\": true, \
        			    \"AttributeName\": \"ttl\" \
    				  } \
                    }' \
            && echo 'Success: TTL enabled for table \"$(SCYLLA_TABLE)\".' \
			|| echo 'Error: TTL enabling request for table \"$(SCYLLA_TABLE)\" failed!';"

disable-ttl: check-scylla-api
	docker exec $(SCYLLA_CONTAINER) /bin/bash -c "curl -o /dev/null -s -f -X POST http://$(SCYLLA_HOST):$(SCYLLA_PORT) \
            	-H 'Content-Type: application/json' \
            	-H 'X-Amz-Target: DynamoDB_20120810.UpdateTimeToLive' \
            	-d '{\
                      \"TableName\": \"$(SCYLLA_TABLE)\", \
    				  \"TimeToLiveSpecification\": { \
        			    \"Enabled\": false, \
        			    \"AttributeName\": \"ttl\" \
    				  } \
                    }' \
            && echo 'Success: TTL disabled for table \"$(SCYLLA_TABLE)\".' \
			|| echo 'Error: TTL disabling request for table \"$(SCYLLA_TABLE)\" failed!';"
