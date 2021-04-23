package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/containerd/containerd/content"
	"github.com/containerd/containerd/content/local"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/docker/distribution"
	"github.com/docker/distribution/reference"
	"github.com/docker/distribution/registry/api/errcode"
	distv2 "github.com/docker/distribution/registry/api/v2"
	"github.com/gorilla/mux"
)

// ResolverProvider provides new resolver
type ResolverProvider func() remotes.Resolver

type Registry struct {
	Resolver ResolverProvider
	Store    content.Store

	srv *http.Server
}

var PLAIN_HTTP_RESOLVER = docker.NewResolver(docker.ResolverOptions{PlainHTTP: true})

func NewRegistry() (*Registry, error) {
	storePath := "/var/facade/localstore"
	store, err := local.NewStore(storePath)
	if err != nil {
		return nil, err
	}

	return &Registry{
		Resolver: func() remotes.Resolver { return PLAIN_HTTP_RESOLVER },
		Store:    store,
	}, nil
}

// MustServe calls serve and logs any error as Fatal
func (reg *Registry) MustServe() {
	err := reg.Serve()
	if err != nil {
		log.Fatal("cannot serve registry")
	}
}

// Serve serves the registry on the given port
func (reg *Registry) Serve() error {
	routes := distv2.RouterWithPrefix("/")
	reg.registerHandler(routes)

	var handler http.Handler = routes
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	addr := ":6000"
	var (
		l   net.Listener
		err error
	)
	// l, err = ReceiveHandover(ctx, reg.Config.Handover.Sockets)
	l, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	reg.srv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	// handover?
	// tls?
	return reg.srv.Serve(l)
}

// registerHandler registers the handle* functions with the corresponding routes
func (reg *Registry) registerHandler(routes *mux.Router) {
	routes.Get(distv2.RouteNameBase).HandlerFunc(reg.handleAPIBase)
	routes.Get(distv2.RouteNameManifest).Handler(dispatcher(reg.handleManifest))
	// routes.Get(v2.RouteNameCatalog).Handler(dispatcher(reg.handleCatalog))
	// routes.Get(v2.RouteNameTags).Handler(dispatcher(reg.handleTags))
	routes.Get(distv2.RouteNameBlob).Handler(dispatcher(reg.handleBlob))
	// routes.Get(v2.RouteNameBlobUpload).Handler(dispatcher(reg.handleBlobUpload))
	// routes.Get(v2.RouteNameBlobUploadChunk).Handler(dispatcher(reg.handleBlobUploadChunk))
	routes.NotFoundHandler = http.HandlerFunc(reg.handleAPIBase)
}

// handleApiBase implements a simple yes-man for doing overall checks against the
// api. This can support auth roundtrips to support docker login.
func (reg *Registry) handleAPIBase(w http.ResponseWriter, r *http.Request) {
	const emptyJSON = "{}"
	// Provide a simple /v2/ 200 OK response with empty json response.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprint(len(emptyJSON)))

	fmt.Fprint(w, emptyJSON)
}

type dispatchFunc func(ctx context.Context, r *http.Request) http.Handler

// dispatcher wraps a dispatchFunc and provides context
func dispatcher(d dispatchFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fc, _ := httputil.DumpRequest(r, false)
		log.Printf("dispatching request: %s\nvars: %v", string(fc), mux.Vars(r))

		// Get context from request, add vars and other info and sync back
		ctx := r.Context()
		ctx = &muxVarsContext{
			Context: ctx,
			vars:    mux.Vars(r),
		}
		r = r.WithContext(ctx)

		if nameRequired(r) {
			nameRef, err := reference.WithName(getName(ctx))
			if err != nil {
				log.Printf("error parsing reference from context – nameRef: %s – err: %w", nameRef, err)
				respondWithError(w, distribution.ErrRepositoryNameInvalid{
					Name:   nameRef.Name(),
					Reason: err,
				})
				return
			}
		}

		d(ctx, r).ServeHTTP(w, r)
	})
}

func respondWithError(w http.ResponseWriter, terr error) {
	err := errcode.ServeJSON(w, terr)
	if err != nil {
		log.Printf("error serving error json – err: %w – origerr: %w", err, terr)
	}
}

// nameRequired returns true if the route requires a name.
func nameRequired(r *http.Request) bool {
	route := mux.CurrentRoute(r)
	if route == nil {
		return true
	}
	routeName := route.GetName()
	return routeName != distv2.RouteNameBase && routeName != distv2.RouteNameCatalog
}

type muxVarsContext struct {
	context.Context
	vars map[string]string
}

func (ctx *muxVarsContext) Value(key interface{}) interface{} {
	if keyStr, ok := key.(string); ok {
		if keyStr == "vars" {
			return ctx.vars
		}

		if strings.HasPrefix(keyStr, "vars.") {
			keyStr = strings.TrimPrefix(keyStr, "vars.")
		}

		if v, ok := ctx.vars[keyStr]; ok {
			return v
		}
	}

	return ctx.Context.Value(key)
}

// getName extracts the name var from the context which was passed in through the mux route
func getName(ctx context.Context) string {
	val := ctx.Value("vars.name")
	sval, ok := val.(string)
	if !ok {
		return ""
	}
	return sval
}

func getSpecProviderName(ctx context.Context) (specProviderName string, remainder string) {
	name := getName(ctx)
	segs := strings.Split(name, "/")
	if len(segs) > 1 {
		specProviderName = segs[0]
		remainder = strings.Join(segs[1:], "/")
	}
	return
}

// getReference extracts the referece var from the context which was passed in through the mux route
func getReference(ctx context.Context) string {
	val := ctx.Value("vars.reference")
	sval, ok := val.(string)
	if !ok {
		return ""
	}
	return sval
}

// getDigest extracts the digest var from the context which was passed in through the mux route
func getDigest(ctx context.Context) string {
	val := ctx.Value("vars.digest")
	sval, ok := val.(string)
	if !ok {
		return ""
	}

	return sval
}
