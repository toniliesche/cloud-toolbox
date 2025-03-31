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

ARG DEBIANVER=bookworm
ARG GOLANGVER=1.24
FROM golang:${GOLANGVER}-${DEBIANVER} AS builder

ARG CACHEBUST
COPY . /tmp/cloud-toolbox

WORKDIR /tmp/cloud-toolbox

ARG BUILDVER
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make build-event-publisher && \
    make build-hello-world-faas

FROM debian:${DEBIANVER}-slim

ARG BUILDVER
ARG COMMIT
LABEL org.opencontainers.image.authors="Toni Liesche <dev@toniliesche.de>"
LABEL org.opencontainers.image.url="https://hub.docker.com/r/tliesche/event-publisher"
LABEL org.opencontainers.image.documentation="https://raw.githubusercontent.com/toniliesche/cloud-toolbox/refs/heads/develop/README.md"
LABEL org.opencontainers.image.source="https://github.com/toniliesche/cloud-toolbox"
LABEL org.opencontainers.image.version="${BUILDVER} <${COMMIT}>"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="Event Publisher container"
LABEL org.opencontainers.image.description="HTTP triggered Event Publisher container"
ENV HTTP_SERVER_HOST=0.0.0.0

COPY --from=builder /tmp/cloud-toolbox/build/event-publisher /usr/local/bin/event-publisher
COPY --from=builder /tmp/cloud-toolbox/build/hello-world-faas /usr/local/bin/hello-world-faas

CMD [ "/usr/local/bin/event-publisher" ]
