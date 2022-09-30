package testdata

var SimpleOriginHandleCode = `// HandleSimpleOriginOrigin applies the CORS response headers corresponding to
// the origin for the service SimpleOrigin.
func HandleSimpleOriginOrigin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "SimpleOrigin") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
`

var RegexpOriginHandleCode = `// HandleRegexpOriginOrigin applies the CORS response headers corresponding to
// the origin for the service RegexpOrigin.
func HandleRegexpOriginOrigin(h http.Handler) http.Handler {
	spec0 := regexp.MustCompile(".*RegexpOrigin.*")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOriginRegexp(origin, spec0) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
`

var MultiOriginHandleCode = `// HandleMultiOriginOrigin applies the CORS response headers corresponding to
// the origin for the service MultiOrigin.
func HandleMultiOriginOrigin(h http.Handler) http.Handler {
	spec0 := regexp.MustCompile(".*MultiOrigin2.*")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOriginRegexp(origin, spec0) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Expose-Headers", "X-Time, X-Api-Version")
			w.Header().Set("Access-Control-Max-Age", "100")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			}
			h.ServeHTTP(w, r)
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
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
`

var OriginFileServerHandleCode = `// HandleOriginFileServerOrigin applies the CORS response headers corresponding
// to the origin for the service OriginFileServer.
func HandleOriginFileServerOrigin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "OriginFileServer") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
`

var OriginMultiEndpointHandleCode = `// HandleOriginMultiEndpointOrigin applies the CORS response headers
// corresponding to the origin for the service OriginMultiEndpoint.
func HandleOriginMultiEndpointOrigin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "OriginMultiEndpoint") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
`

var MultiServiceSameOriginFirstServiceHandleCode = `// HandleFirstServiceOrigin applies the CORS response headers corresponding to
// the origin for the service FirstService.
func HandleFirstServiceOrigin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "SimpleOrigin") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
`
var MultiServiceSameOriginSecondServiceHandleCode = `// HandleSecondServiceOrigin applies the CORS response headers corresponding to
// the origin for the service SecondService.
func HandleSecondServiceOrigin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "SimpleOrigin") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
`

var FilesHandleCode = `// HandleFilesOrigin applies the CORS response headers corresponding to the
// origin for the service Files.
func HandleFilesOrigin(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			// Not a CORS request
			h.ServeHTTP(w, r)
			return
		}
		if cors.MatchOrigin(origin, "*") {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			if acrm := r.Header.Get("Access-Control-Request-Method"); acrm != "" {
				// We are handling a preflight request
			}
			h.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
		return
	})
}
`

var SimpleOriginMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service SimpleOrigin.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleSimpleOriginOrigin(h)
	mux.Handle("OPTIONS", "/", h.ServeHTTP)
}
`

var RegexpOriginMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service RegexpOrigin.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleRegexpOriginOrigin(h)
	mux.Handle("OPTIONS", "/", h.ServeHTTP)
}
`

var MultiOriginMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service MultiOrigin.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleMultiOriginOrigin(h)
	mux.Handle("OPTIONS", "/", h.ServeHTTP)
}
`

var OriginFileServerMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service OriginFileServer.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleOriginFileServerOrigin(h)
	mux.Handle("OPTIONS", "/file.json", h.ServeHTTP)
}
`

var OriginMultiEndpointMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service OriginMultiEndpoint.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleOriginMultiEndpointOrigin(h)
	mux.Handle("OPTIONS", "/{:id}", h.ServeHTTP)
	mux.Handle("OPTIONS", "/", h.ServeHTTP)
}
`

var MultiServiceSameOriginFirstServiceMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service FirstService.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleFirstServiceOrigin(h)
	mux.Handle("OPTIONS", "/", h.ServeHTTP)
}
`
var MultiServiceSameOriginSecondServiceMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service SecondService.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleSecondServiceOrigin(h)
	mux.Handle("OPTIONS", "/", h.ServeHTTP)
}
`

var FilesMountCode = `// MountCORSHandler configures the mux to serve the CORS endpoints for the
// service Files.
func MountCORSHandler(mux goahttp.Muxer, h http.Handler) {
	h = HandleFilesOrigin(h)
	mux.Handle("OPTIONS", "/index", h.ServeHTTP)
}
`

var SimpleOriginServerInitCode = `// New instantiates HTTP handlers for all the SimpleOrigin service endpoints
// using the provided encoder and decoder. The handlers are mounted on the
// given mux using the HTTP verb and path defined in the design. errhandler is
// called whenever a response fails to be encoded. formatter is used to format
// errors returned by the service methods prior to encoding. Both errhandler
// and formatter are optional and can be nil.
func New(
	e *simpleorigin.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"SimpleOriginMethod", "GET", "/"},
			{"CORS", "OPTIONS", "/"},
		},
		SimpleOriginMethod: NewSimpleOriginMethodHandler(e.SimpleOriginMethod, mux, decoder, encoder, errhandler, formatter),
		CORS:               NewCORSHandler(),
	}
}
`

var RegexpOriginServerInitCode = `// New instantiates HTTP handlers for all the RegexpOrigin service endpoints
// using the provided encoder and decoder. The handlers are mounted on the
// given mux using the HTTP verb and path defined in the design. errhandler is
// called whenever a response fails to be encoded. formatter is used to format
// errors returned by the service methods prior to encoding. Both errhandler
// and formatter are optional and can be nil.
func New(
	e *regexporigin.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"RegexpOriginMethod", "GET", "/"},
			{"CORS", "OPTIONS", "/"},
		},
		RegexpOriginMethod: NewRegexpOriginMethodHandler(e.RegexpOriginMethod, mux, decoder, encoder, errhandler, formatter),
		CORS:               NewCORSHandler(),
	}
}
`

var MultiOriginServerInitCode = `// New instantiates HTTP handlers for all the MultiOrigin service endpoints
// using the provided encoder and decoder. The handlers are mounted on the
// given mux using the HTTP verb and path defined in the design. errhandler is
// called whenever a response fails to be encoded. formatter is used to format
// errors returned by the service methods prior to encoding. Both errhandler
// and formatter are optional and can be nil.
func New(
	e *multiorigin.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"MultiOriginMethod", "GET", "/"},
			{"CORS", "OPTIONS", "/"},
		},
		MultiOriginMethod: NewMultiOriginMethodHandler(e.MultiOriginMethod, mux, decoder, encoder, errhandler, formatter),
		CORS:              NewCORSHandler(),
	}
}
`

var OriginFileServerServerInitCode = `// New instantiates HTTP handlers for all the OriginFileServer service
// endpoints using the provided encoder and decoder. The handlers are mounted
// on the given mux using the HTTP verb and path defined in the design.
// errhandler is called whenever a response fails to be encoded. formatter is
// used to format errors returned by the service methods prior to encoding.
// Both errhandler and formatter are optional and can be nil.
func New(
	e *originfileserver.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
	fileSystemFileJSON http.FileSystem,
) *Server {
	if fileSystemFileJSON == nil {
		fileSystemFileJSON = http.Dir(".")
	}
	return &Server{
		Mounts: []*MountPoint{
			{"CORS", "OPTIONS", "/file.json"},
			{"./file.json", "GET", "/file.json"},
		},
		CORS:     NewCORSHandler(),
		FileJSON: http.FileServer(fileSystemFileJSON),
	}
}
`

var OriginMultiEndpointServerInitCode = `// New instantiates HTTP handlers for all the OriginMultiEndpoint service
// endpoints using the provided encoder and decoder. The handlers are mounted
// on the given mux using the HTTP verb and path defined in the design.
// errhandler is called whenever a response fails to be encoded. formatter is
// used to format errors returned by the service methods prior to encoding.
// Both errhandler and formatter are optional and can be nil.
func New(
	e *originmultiendpoint.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"OriginMultiEndpointGet", "GET", "/{:id}"},
			{"OriginMultiEndpointPost", "POST", "/"},
			{"OriginMultiEndpointOptions", "OPTIONS", "/ids/{:id}"},
			{"CORS", "OPTIONS", "/{:id}"},
			{"CORS", "OPTIONS", "/"},
		},
		OriginMultiEndpointGet:     NewOriginMultiEndpointGetHandler(e.OriginMultiEndpointGet, mux, decoder, encoder, errhandler, formatter),
		OriginMultiEndpointPost:    NewOriginMultiEndpointPostHandler(e.OriginMultiEndpointPost, mux, decoder, encoder, errhandler, formatter),
		OriginMultiEndpointOptions: NewOriginMultiEndpointOptionsHandler(e.OriginMultiEndpointOptions, mux, decoder, encoder, errhandler, formatter),
		CORS:                       NewCORSHandler(),
	}
}
`

var MultiServiceSameOriginFirstServiceInitCode = `// New instantiates HTTP handlers for all the FirstService service endpoints
// using the provided encoder and decoder. The handlers are mounted on the
// given mux using the HTTP verb and path defined in the design. errhandler is
// called whenever a response fails to be encoded. formatter is used to format
// errors returned by the service methods prior to encoding. Both errhandler
// and formatter are optional and can be nil.
func New(
	e *firstservice.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"SimpleOriginMethod", "GET", "/"},
			{"CORS", "OPTIONS", "/"},
		},
		SimpleOriginMethod: NewSimpleOriginMethodHandler(e.SimpleOriginMethod, mux, decoder, encoder, errhandler, formatter),
		CORS:               NewCORSHandler(),
	}
}
`

var MultiServiceSameOriginSecondServiceInitCode = `// New instantiates HTTP handlers for all the SecondService service endpoints
// using the provided encoder and decoder. The handlers are mounted on the
// given mux using the HTTP verb and path defined in the design. errhandler is
// called whenever a response fails to be encoded. formatter is used to format
// errors returned by the service methods prior to encoding. Both errhandler
// and formatter are optional and can be nil.
func New(
	e *secondservice.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"SimpleOriginMethod", "GET", "/"},
			{"CORS", "OPTIONS", "/"},
		},
		SimpleOriginMethod: NewSimpleOriginMethodHandler(e.SimpleOriginMethod, mux, decoder, encoder, errhandler, formatter),
		CORS:               NewCORSHandler(),
	}
}
`

var FilesServerInitCode = `// New instantiates HTTP handlers for all the Files service endpoints using the
// provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *files.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(ctx context.Context, err error) goahttp.Statuser,
	fileSystemIndexHTML http.FileSystem,
) *Server {
	if fileSystemIndexHTML == nil {
		fileSystemIndexHTML = http.Dir(".")
	}
	return &Server{
		Mounts: []*MountPoint{
			{"CORS", "OPTIONS", "/index"},
			{"index.html", "GET", "/index"},
		},
		CORS:      NewCORSHandler(),
		IndexHTML: http.FileServer(fileSystemIndexHTML),
	}
}
`
