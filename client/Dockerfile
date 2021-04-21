FROM golang

WORKDIR $GOPATH/src/gitpod.io/containerd-test-client
COPY go.mod go.sum main.go ./
RUN go get -d -v ./...
RUN go install -v ./...
ENTRYPOINT [ "containerd-test-client" ]
