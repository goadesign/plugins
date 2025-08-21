{{ printf "Harness manages test servers and provides a test client for the %s service." .Name | comment }}
type Harness struct {
	t         *testing.T
	service   {{ .PkgName }}.Service
	Client    *Client
	ctx       context.Context
{{- if .HasHTTP }}
	httpSvr   *httptest.Server
	httpCli   *http.Client
{{- end }}
{{- if .HasGRPC }}
	grpcConn  *grpc.ClientConn
	grpcSvr   *grpc.Server
	grpcLis   *bufconn.Listener
{{- end }}
{{- if .HasJSONRPC }}
	jsonrpcSvr *httptest.Server
	jsonrpcCli *http.Client
{{- end }}
}

{{ printf "Options configure the test harness." | comment }}
type Options struct {
	Context context.Context
	Timeout time.Duration
}

{{ printf "Option is a harness configuration option." | comment }}
type Option func(*Options)

{{ printf "WithContext sets a custom context for the harness." | comment }}
func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

{{ printf "WithTimeout sets a timeout for all operations." | comment }}
func WithTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.Timeout = d
	}
}
