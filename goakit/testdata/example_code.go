package testdata

var MultiEndpointServiceStructCode = `// MultiEndpointService service example implementation.
// The example methods log the requests and return zero values.
type multiendpointservicesvc struct {
	logger log.Logger
}

// NewMultiEndpointService returns the MultiEndpointService service
// implementation.
func NewMultiEndpointService(logger log.Logger) multiendpointservice.Service {
	return &multiendpointservicesvc{logger}
}
`

var MixedServiceStructCode = `// MixedService service example implementation.
// The example methods log the requests and return zero values.
type mixedservicesvc struct {
	logger log.Logger
}

// NewMixedService returns the MixedService service implementation.
func NewMixedService(logger log.Logger) mixedservice.Service {
	return &mixedservicesvc{logger}
}
`
