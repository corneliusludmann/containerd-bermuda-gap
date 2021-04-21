# FROM alpine
# RUN apk add --no-cache \
#     curl file

FROM ubuntu
RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

ARG VERSION=1.2.10
ARG DELIM=.
RUN curl -sSOL https://github.com/containerd/containerd/releases/download/v${VERSION}/containerd-${VERSION}${DELIM}linux-amd64.tar.gz && \
    tar -xvzf containerd-${VERSION}${DELIM}linux-amd64.tar.gz

# COPY containerd-config.toml /etc/containerd/config.toml

ENTRYPOINT [ "containerd" ]
