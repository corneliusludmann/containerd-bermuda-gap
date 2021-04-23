package main

import (
	"context"
	"log"
	"net/http"

	distv2 "github.com/docker/distribution/registry/api/v2"
)

func (reg *Registry) handleBlob(ctx context.Context, r *http.Request) http.Handler {
	spname, name := getSpecProviderName(ctx)
	log.Printf("handleBlob – unknown spec provider – spname: %s – name: %s", spname, name)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, distv2.ErrorCodeManifestUnknown)
	})
}
