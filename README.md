# containerd-bermuda-gap

## Push images to the registry

```
docker-compose up -d registry
docker pull gitpod/workspace-full
docker tag gitpod/workspace-full localhost:5000/workspace-full:latest
docker push localhost:5000/workspace-full:latest
docker-compose down
```

## Vagrant

```
$ vagrant plugin install vagrant-disksize
$ mkdir -p logs
$ vagrant destroy -f && vagrant up
$ vagrant snapshot save init --force    # optional
$ vagrant ssh
$ sudo /root/go/bin/containerd-test-client localhost:5000/workspace-full:latest 10 2>&1 | tee /home/vagrant/logs/client.log
```

Facade:
```
$ vagrant ssh
$ sudo su -
$ cp -r /home/vagrant/facade /root/go/src/gitpod.io/
$ cd /root/go/src/gitpod.io/facade
$ export PATH=$PATH:/usr/local/go/bin
$ export GOPATH=/root/go
$ go get -v ./...
$ go install -v ./...
$ /root/go/bin/facade 2>&1 | tee /home/vagrant/logs/facade.log
$ /root/go/bin/containerd-test-client localhost:6000/workspace-full:latest 10 2>&1 | tee /home/vagrant/logs/client.log
```
