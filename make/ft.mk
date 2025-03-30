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

setup-ft:
	$(call print_message,Setting up RabbitMQ for FaaS)
	$(MAKE) setup-ft-rabbitmq

setup-ft-rabbitmq:
	$(call print_message,Creating RabbitMQ User "$(FT_RABBITMQ_USER)")
	$(MAKE) create-rabbitmq-user RABBITMQ_USER=$(FT_RABBITMQ_USER) RABBITMQ_PASSWORD=$(FT_RABBITMQ_PASSWORD)

	$(call print_message,Creating RabbitMQ vHost "$(FT_RABBITMQ_VHOST)")
	$(MAKE) create-rabbitmq-vhost RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST)

	$(call print_message,Setting RabbitMQ Permissions for User "$(RABBITMQ_ADMIN_USER)" on vHost "$(FT_RABBITMQ_VHOST)")
	$(MAKE) set-rabbitmq-user-permissions RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST) RABBITMQ_USER=$(RABBITMQ_ADMIN_USER) RABBITMQ_PERM_CONF=".*" RABBITMQ_PERM_WRITE=".*" RABBITMQ_PERM_READ=".*"

	$(call print_message,Setting RabbitMQ Permissions for User "$(FT_RABBITMQ_USER)" on vHost "$(FT_RABBITMQ_VHOST)")
	$(MAKE) set-rabbitmq-user-permissions RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST) RABBITMQ_USER=$(FT_RABBITMQ_USER) RABBITMQ_PERM_CONF="" RABBITMQ_PERM_WRITE="" RABBITMQ_PERM_READ="ft.*"

	$(call print_message,Creating RabbitMQ Exchange "$(FT_RABBITMQ_DLX)")
	$(MAKE) create-rabbitmq-exchange RABBITMQ_EXCHANGE=$(FT_RABBITMQ_DLX) RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST)

	$(call print_message,Creating RabbitMQ Queue "$(FT_RABBITMQ_QUEUE).dlq")
	$(MAKE) create-rabbitmq-queue RABBITMQ_QUEUE=$(FT_RABBITMQ_QUEUE).dlq RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST)

	$(call print_message,Creating RabbitMQ Queue "$(FT_RABBITMQ_QUEUE)")
	$(MAKE) create-rabbitmq-queue RABBITMQ_QUEUE=$(FT_RABBITMQ_QUEUE) RABBITMQ_DLX=$(FT_RABBITMQ_DLX) RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST)

	$(call print_message,Creating RabbitMQ Binding from Exchange "$(FT_RABBITMQ_DLX)" to Queue "$(FT_RABBITMQ_QUEUE).dlq" (Topic: "$(FT_RABBITMQ_QUEUE).dlq"))
	$(MAKE) create-rabbitmq-binding RABBITMQ_EXCHANGE=$(FT_RABBITMQ_DLX) RABBITMQ_QUEUE=$(FT_RABBITMQ_QUEUE).dlq RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST) RABBITMQ_BINDING_KEY="$(FT_RABBITMQ_QUEUE).dlq"

	$(call print_message,Creating RabbitMQ Exchange "$(FT_RABBITMQ_EXCHANGE)")
	$(MAKE) create-rabbitmq-exchange RABBITMQ_EXCHANGE=$(FT_RABBITMQ_EXCHANGE) RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST)

	$(call print_message,Creating RabbitMQ Binding from Exchange "$(FT_RABBITMQ_EXCHANGE)" to Queue "$(FT_RABBITMQ_QUEUE)" (Topic: "$(FT_RABBITMQ_QUEUE)"))
	$(MAKE) create-rabbitmq-binding RABBITMQ_EXCHANGE=$(FT_RABBITMQ_EXCHANGE) RABBITMQ_QUEUE=$(FT_RABBITMQ_QUEUE) RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST) RABBITMQ_BINDING_KEY="$(FT_RABBITMQ_QUEUE)"

uninstall-ft:
	$(call print_message,Deleting RabbitMQ vHost "$(FT_RABBITMQ_VHOST)")
	$(MAKE) delete-rabbitmq-vhost RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST)

	$(call print_message,Deleting RabbitMQ User "$(FT_RABBITMQ_USER)")
	$(MAKE) delete-rabbitmq-user RABBITMQ_USER=$(FT_RABBITMQ_USER)