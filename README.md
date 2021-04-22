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
$ vagrant ssh
$ /root/go/bin/containerd-test-client localhost:5000/workspace-full:latest 10 2>&1 | tee /home/vagrant/logs/client.log
```