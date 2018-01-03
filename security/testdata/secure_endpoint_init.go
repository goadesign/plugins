package testdata

var EndpointInitWithoutRequirementCode = `// NewSecureEndpoints wraps the methods of a EndpointWithoutRequirement service
// with security scheme aware endpoints.
func NewSecureEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Unsecure: NewUnsecure(s),
	}
}
`

var EndpointInitWithRequirementsCode = `// NewSecureEndpoints wraps the methods of a EndpointsWithRequirements service
// with security scheme aware endpoints.
func NewSecureEndpoints(s Service) *Endpoints {
	return &Endpoints{
		SecureWithRequirements:       SecureSecureWithRequirements(NewSecureWithRequirementsEndpoint(s)),
		DoublySecureWithRequirements: SecureDoublySecureWithRequirements(NewDoublySecureWithRequirementsEndpoint(s)),
	}
}
`

var EndpointInitWithServiceRequirementsCode = `// NewSecureEndpoints wraps the methods of a EndpointsWithServiceRequirements
// service with security scheme aware endpoints.
func NewSecureEndpoints(s Service) *Endpoints {
	return &Endpoints{
		SecureWithRequirements:     SecureSecureWithRequirements(NewSecureWithRequirementsEndpoint(s)),
		AlsoSecureWithRequirements: SecureAlsoSecureWithRequirements(NewAlsoSecureWithRequirementsEndpoint(s)),
	}
}
`

var EndpointNoSecurityCode = `// NewSecureEndpoints wraps the methods of a EndpointNoSecurity service with
// security scheme aware endpoints.
func NewSecureEndpoints(s Service) *Endpoints {
	return &Endpoints{
		NoSecurity: SecureNoSecurity(NewNoSecurityEndpoint(s)),
	}
}
`
