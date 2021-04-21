package main

import (
	"context"
	"log"
	"os"
	"strings"

	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/remotes/docker"
)

var DEFAULT_REF = "docker.io/library/alpine:latest"
var PLAIN_HTTP_RESOLVER = docker.NewResolver(docker.ResolverOptions{PlainHTTP: true})

func main() {
	log.Println("Sleeping ...")
	time.Sleep(10 * time.Second)
	ref := DEFAULT_REF
	if len(os.Args) > 1 {
		ref = os.Args[1]
	}
	log.Printf("Starting using image '%s' ...\n", ref)

	for {
		if err := pullImage(context.Background(), ref); err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Second)
	}
}

func pullImage(ctx context.Context, ref string) error {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	ctx2 := namespaces.WithNamespace(ctx, "client")

	var image containerd.Image
	if strings.HasPrefix(ref, "registry:5000") {
		image, err = client.Pull(ctx2, ref, containerd.WithPullUnpack, containerd.WithResolver(PLAIN_HTTP_RESOLVER))
	} else {
		image, err = client.Pull(ctx2, ref, containerd.WithPullUnpack)
	}
	if err != nil {
		return err
	}
	log.Printf("Successfully pulled '%s' image.\n", image.Name())

	return nil
}
