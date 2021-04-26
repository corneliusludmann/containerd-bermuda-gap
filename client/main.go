package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/remotes/docker"
)

var DEFAULT_REF = "docker.io/library/alpine:latest"
var PLAIN_HTTP_RESOLVER = docker.NewResolver(docker.ResolverOptions{PlainHTTP: true})

func main() {
	// log.Println("Sleeping ...")
	// time.Sleep(10 * time.Second)
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

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var wg sync.WaitGroup

	i := 0
	for {
		wg.Add(1)
		go func(i int) {

			defer wg.Done()
			if max > 1 && i > (max/2) {
				delay := 30 * time.Second
				log.Printf("Sleeping %s until pulling %d ...\n", delay, i)
				time.Sleep(delay)
			}
			log.Printf("Pulling %d ...\n", i)
			if err := pullImage(context.Background(), client, ref, i); err != nil {
				log.Fatal(err)
			}
		}(i)
		i++
		if i == max {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	wg.Wait()
	log.Println("Done.")
}

func pullImage(ctx context.Context, client *containerd.Client, ref string, i int) error {
	defer timeTrack(time.Now(), fmt.Sprintf("Pulling %d", i))

	ctx2 := namespaces.WithNamespace(ctx, "client")
	var err error = nil

	var image containerd.Image
	if strings.HasPrefix(ref, "registry:5000") {
		// without native snapshotter I'm able to pull alpine but not workspace-full (in Docker)
		// error: failed to extract layer [...] failed to mount /var/lib/containerd/tmpmounts/containerd-mount366656299: invalid argument: unknown
		// see also: https://github.com/containerd/containerd/issues/2402#issuecomment-398033418
		//image, err = client.Pull(ctx2, ref, containerd.WithPullUnpack, containerd.WithPullSnapshotter("native"), containerd.WithResolver(PLAIN_HTTP_RESOLVER))
		image, err = client.Pull(ctx2, ref, containerd.WithPullUnpack, containerd.WithResolver(PLAIN_HTTP_RESOLVER))
	} else {
		//image, err = client.Pull(ctx2, ref, containerd.WithPullUnpack, containerd.WithPullSnapshotter("native"))
		image, err = client.Pull(ctx2, ref, containerd.WithPullUnpack)
	}
	if err != nil {
		return err
	}
	log.Printf("Successfully pulled '%s' image (%d).\n", image.Name(), i)

	return nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s.", name, elapsed)
}
