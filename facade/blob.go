package main

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (reg *Registry) handleBlob(ctx context.Context, r *http.Request) http.Handler {
	spname, name := getSpecProviderName(ctx)
	reference := getReference(ctx)
	log.Printf("handleBlob – spname: %s – name: %s – reference: %s", spname, name, reference)
	target, _ := url.Parse("http://localhost:5000/")
	proxy := httputil.NewSingleHostReverseProxy(target)
	return proxy
}
