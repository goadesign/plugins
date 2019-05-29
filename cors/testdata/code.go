package testdata

var SimpleOriginHandleCode = `// handleSimpleOriginOrigin applies the CORS response headers corresponding to
// the origin for the service SimpleOrigin.
func handleSimpleOriginOrigin(h http.Handler) http.Handler {
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "SimpleOrigin") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
`

var RegexpOriginHandleCode = `// handleRegexpOriginOrigin applies the CORS response headers corresponding to
// the origin for the service RegexpOrigin.
func handleRegexpOriginOrigin(h http.Handler) http.Handler {
	spec0 := regexp.MustCompile(".*RegexpOrigin.*")
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOriginRegexp(origin, spec0) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
`

var MultiOriginHandleCode = `// handleMultiOriginOrigin applies the CORS response headers corresponding to
// the origin for the service MultiOrigin.
func handleMultiOriginOrigin(h http.Handler) http.Handler {
	spec0 := regexp.MustCompile(".*MultiOrigin2.*")
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOriginRegexp(origin, spec0) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "X-Time, X-Api-Version")
			w.Header().Set("Access-Control-Max-Age", "100")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			}
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "MultiOrigin1") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "X-Time")
			w.Header().Set("Access-Control-Max-Age", "600")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
				w.Header().Set("Access-Control-Allow-Headers", "X-Shared-Secret")
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
`

var OriginFileServerHandleCode = `// handleOriginFileServerOrigin applies the CORS response headers corresponding
// to the origin for the service OriginFileServer.
func handleOriginFileServerOrigin(h http.Handler) http.Handler {
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "OriginFileServer") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
`

var OriginMultiEndpointHandleCode = `// handleOriginMultiEndpointOrigin applies the CORS response headers
// corresponding to the origin for the service OriginMultiEndpoint.
func handleOriginMultiEndpointOrigin(h http.Handler) http.Handler {
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "OriginMultiEndpoint") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
`

var MultiServiceSameOriginFirstServiceHandleCode = `// handleFirstServiceOrigin applies the CORS response headers corresponding to
// the origin for the service FirstService.
func handleFirstServiceOrigin(h http.Handler) http.Handler {
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "SimpleOrigin") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
`
var MultiServiceSameOriginSecondServiceHandleCode = `// handleSecondServiceOrigin applies the CORS response headers corresponding to
// the origin for the service SecondService.
func handleSecondServiceOrigin(h http.Handler) http.Handler {
	origHndlr := h.(http.HandlerFunc)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			origHndlr(w, r)
			return
		}
		if cors.MatchOrigin(origin, "SimpleOrigin") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "false")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			origHndlr(w, r)
			return
		}
		origHndlr(w, r)
		return
	})
}
`

var SimpleOriginMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service SimpleOrigin.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleSimpleOriginOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/", f)
}
`

var RegexpOriginMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service RegexpOrigin.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleRegexpOriginOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/", f)
}
`

var MultiOriginMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service MultiOrigin.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleMultiOriginOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/", f)
}
`

var OriginFileServerMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service OriginFileServer.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleOriginFileServerOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/file.json", f)
}
`

var OriginMultiEndpointMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service OriginMultiEndpoint.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleOriginMultiEndpointOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/{:id}", f)
	mux.Handle("OPTIONS", "/", f)
}
`

var MultiServiceSameOriginFirstServiceMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service FirstService.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleFirstServiceOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/", f)
}
`
var MultiServiceSameOriginSecondServiceMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service SecondService.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = handleSecondServiceOrigin(h)
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("OPTIONS", "/", f)
}
`

var SimpleOriginServerInitCode = `// New instantiates HTTP handlers for all the SimpleOrigin service endpoints.
func New(
	e *simpleorigin.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
	eh func(context.Context, http.ResponseWriter, error),
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"SimpleOriginMethod", "GET", "/"},
			{"CORS", "OPTIONS", "/"},
		},
		SimpleOriginMethod: NewSimpleOriginMethodHandler(e.SimpleOriginMethod, mux, dec, enc, eh),
		CORS:               NewCORSHandler(),
	}
}
`

var RegexpOriginServerInitCode = `// New instantiates HTTP handlers for all the RegexpOrigin service endpoints.
func New(
	e *regexporigin.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
	eh func(context.Context, http.ResponseWriter, error),
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"RegexpOriginMethod", "GET", "/"},
			{"CORS", "OPTIONS", "/"},
		},
		RegexpOriginMethod: NewRegexpOriginMethodHandler(e.RegexpOriginMethod, mux, dec, enc, eh),
		CORS:               NewCORSHandler(),
	}
}
`

var MultiOriginServerInitCode = `// New instantiates HTTP handlers for all the MultiOrigin service endpoints.
func New(
	e *multiorigin.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
	eh func(context.Context, http.ResponseWriter, error),
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"MultiOriginMethod", "GET", "/"},
			{"CORS", "OPTIONS", "/"},
		},
		MultiOriginMethod: NewMultiOriginMethodHandler(e.MultiOriginMethod, mux, dec, enc, eh),
		CORS:              NewCORSHandler(),
	}
}
`

var OriginFileServerServerInitCode = `// New instantiates HTTP handlers for all the OriginFileServer service
// endpoints.
func New(
	e *originfileserver.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
	eh func(context.Context, http.ResponseWriter, error),
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"CORS", "OPTIONS", "/file.json"},
			{"./file.json", "GET", "/file.json"},
		},
		CORS: NewCORSHandler(),
	}
}
`

var OriginMultiEndpointServerInitCode = `// New instantiates HTTP handlers for all the OriginMultiEndpoint service
// endpoints.
func New(
	e *originmultiendpoint.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
	eh func(context.Context, http.ResponseWriter, error),
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"OriginMultiEndpointGet", "GET", "/{:id}"},
			{"OriginMultiEndpointPost", "POST", "/"},
			{"OriginMultiEndpointOptions", "OPTIONS", "/ids/{:id}"},
			{"CORS", "OPTIONS", "/{:id}"},
			{"CORS", "OPTIONS", "/"},
		},
		OriginMultiEndpointGet:     NewOriginMultiEndpointGetHandler(e.OriginMultiEndpointGet, mux, dec, enc, eh),
		OriginMultiEndpointPost:    NewOriginMultiEndpointPostHandler(e.OriginMultiEndpointPost, mux, dec, enc, eh),
		OriginMultiEndpointOptions: NewOriginMultiEndpointOptionsHandler(e.OriginMultiEndpointOptions, mux, dec, enc, eh),
		CORS:                       NewCORSHandler(),
	}
}
`

var MultiServiceSameOriginFirstServiceInitCode = `// New instantiates HTTP handlers for all the FirstService service endpoints.
func New(
	e *firstservice.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
	eh func(context.Context, http.ResponseWriter, error),
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"SimpleOriginMethod", "GET", "/"},
			{"CORS", "OPTIONS", "/"},
		},
		SimpleOriginMethod: NewSimpleOriginMethodHandler(e.SimpleOriginMethod, mux, dec, enc, eh),
		CORS:               NewCORSHandler(),
	}
}
`
var MultiServiceSameOriginSecondServiceInitCode = `// New instantiates HTTP handlers for all the SecondService service endpoints.
func New(
	e *secondservice.Endpoints,
	mux goahttp.Muxer,
	dec func(*http.Request) goahttp.Decoder,
	enc func(context.Context, http.ResponseWriter) goahttp.Encoder,
	eh func(context.Context, http.ResponseWriter, error),
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"SimpleOriginMethod", "GET", "/"},
			{"CORS", "OPTIONS", "/"},
		},
		SimpleOriginMethod: NewSimpleOriginMethodHandler(e.SimpleOriginMethod, mux, dec, enc, eh),
		CORS:               NewCORSHandler(),
	}
}
`
