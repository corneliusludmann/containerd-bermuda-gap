# containerd-bermuda-gap

## Push images to the registry

```
docker-compose up -d registry
docker pull gitpod/workspace-full
docker tag gitpod/workspace-full localhost:5000/workspace-full:latest
docker push localhost:5000/workspace-full:latest
docker-compose down
```
