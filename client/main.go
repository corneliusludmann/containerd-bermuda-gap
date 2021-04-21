package main

import (
	"context"
	"log"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
)

func main() {
	if err := pullImage("docker.io/library/alpine:latest"); err != nil {
		log.Fatal(err)
	}
}

func pullImage(ref string) error {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()
	ctx := namespaces.WithNamespace(context.Background(), "client")
	image, err := client.Pull(ctx, ref, containerd.WithPullUnpack)
	if err != nil {
		return err
	}
	log.Printf("Successfully pulled %s image\n", image.Name())

	return nil
}
