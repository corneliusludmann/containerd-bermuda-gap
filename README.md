# containerd-bermuda-gap

## Push images to the registry

```
docker-compose up -d registry
docker pull alpine
docker tag alpine localhost:5000/alpine:latest
docker push localhost:5000/alpine:latest
docker-compose down
```
