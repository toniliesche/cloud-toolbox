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

setup-ep:
	$(call print_message,Setting up RabbitMQ for EP)
	$(MAKE) setup-ep-rabbitmq

setup-ep-rabbitmq:
	$(call print_message,Creating RabbitMQ User "$(EP_RABBITMQ_USER)")
	$(MAKE) create-rabbitmq-user RABBITMQ_USER=$(EP_RABBITMQ_USER) RABBITMQ_PASSWORD=$(EP_RABBITMQ_PASSWORD)

	$(call print_message,Setting RabbitMQ Permissions for User "$(EP_RABBITMQ_USER)" on vHost "$(FT_RABBITMQ_VHOST)")
	$(MAKE) set-rabbitmq-user-permissions RABBITMQ_VHOST=$(FT_RABBITMQ_VHOST) RABBITMQ_USER=$(EP_RABBITMQ_USER) RABBITMQ_PERM_CONF="" RABBITMQ_PERM_WRITE="${FT_RABBITMQ_EXCHANGE}" RABBITMQ_PERM_READ=""

uninstall-ep:
	$(call print_message,Deleting RabbitMQ User "$(EP_RABBITMQ_USER)")
	$(MAKE) delete-rabbitmq-user RABBITMQ_USER=$(EP_RABBITMQ_USER)