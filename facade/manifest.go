package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func (reg *Registry) handleManifest(ctx context.Context, r *http.Request) http.Handler {
	spname, name := getSpecProviderName(ctx)
	reference := getReference(ctx)
	log.Printf("handleManifest – spname: %s – name: %s – reference: %s", spname, name, reference)
	time.Sleep(10 * time.Second)
	log.Printf("handleManifest – end of sleep – spname: %s – name: %s – reference: %s", spname, name, reference)
	target, _ := url.Parse("http://localhost:5000/")
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}
