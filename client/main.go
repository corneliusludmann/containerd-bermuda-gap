package main

import (
	"context"
	"log"
	"os"
	"strconv"
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
	max := -1
	if len(os.Args) > 2 {
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		max = i
	}
	log.Printf("Starting using image '%s', max='%d' ...\n", ref, max)

	i := 0
	for {
		if err := pullImage(context.Background(), ref); err != nil {
			log.Fatal(err)
		}
		i++
		if i == max {
			break
		}
		time.Sleep(1 * time.Second)
	}
	log.Println("Done.")
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
		// without native snapshotter I'm able to pull alpine but not workspace-full
		// error: failed to extract layer [...] failed to mount /var/lib/containerd/tmpmounts/containerd-mount366656299: invalid argument: unknown
		// see also: https://github.com/containerd/containerd/issues/2402#issuecomment-398033418
		image, err = client.Pull(ctx2, ref, containerd.WithPullUnpack, containerd.WithPullSnapshotter("native"), containerd.WithResolver(PLAIN_HTTP_RESOLVER))
	} else {
		image, err = client.Pull(ctx2, ref, containerd.WithPullUnpack, containerd.WithPullSnapshotter("native"))
	}
	if err != nil {
		return err
	}
	log.Printf("Successfully pulled '%s' image.\n", image.Name())

	return nil
}
