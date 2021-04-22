#!/usr/bin/env bash

set -exuo pipefail

cd /

apt update -y
apt install -y \
    gcc
#    #golang-go

mkdir -p /root/go/src/gitpod.io
mv /home/vagrant/workdir/client /root/go/src/gitpod.io/
mv /home/vagrant/registry /var/lib/registry
mkdir -p /etc/docker/registry

# golang
curl -sSOL https://golang.org/dl/go1.16.3.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/root/go
go version

# containerd
VERSION=1.2.10
DELIM=.
curl -sSOL https://github.com/containerd/containerd/releases/download/v${VERSION}/containerd-${VERSION}${DELIM}linux-amd64.tar.gz && \
tar -xvzf containerd-${VERSION}${DELIM}linux-amd64.tar.gz

# registry
curl -sSLo /etc/docker/registry/config.yml https://raw.githubusercontent.com/docker/distribution-library-image/ab00e8dae12d4515ed259015eab771ec92e92dd4/amd64/config-example.yml
curl -sSLo /usr/local/bin/registry https://raw.githubusercontent.com/docker/distribution-library-image/ab00e8dae12d4515ed259015eab771ec92e92dd4/amd64/registry
chmod +x /usr/local/bin/registry

# start containerd and registry
containerd >/home/vagrant/logs/containerd.log 2>&1 &
registry serve /etc/docker/registry/config.yml >/home/vagrant/logs/registry.log 2>&1 &

# build client
cd /root/go/src/gitpod.io/client
go get -v ./...
go install -v ./...
export PATH=$PATH:/root/go/bin
# containerd-test-client localhost:5000/workspace-full:latest 10 2>&1 | tee /home/vagrant/logs/client.log
