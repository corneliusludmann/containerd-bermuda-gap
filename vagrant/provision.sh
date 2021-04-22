#!/usr/bin/env bash

set -exuo pipefail

VERSION=1.2.10
DELIM=.
cd /
curl -sSOL https://github.com/containerd/containerd/releases/download/v${VERSION}/containerd-${VERSION}${DELIM}linux-amd64.tar.gz && \
tar -xvzf containerd-${VERSION}${DELIM}linux-amd64.tar.gz

containerd --version
