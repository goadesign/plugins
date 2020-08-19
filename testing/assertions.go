package testing

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	goa "goa.design/goa/v3/pkg"
)

type (
	/*// Asserter ensures all test assertions defined on a Goa endpoint
	// succeeds.
	Asserter interface {
		// Met returns an error if there are any assertions that are not met.
		Met() error
	}

	// Mocker mocks a Goa endpoint.
	Mocker interface {
		// Gets mocks the payload received by the Goa endpoint.
		Gets(interface{}) *Endpoint
		// Returns mocks the result returned by the Goa endpoint.
		Returns(interface{}) *Endpoint
		// ReturnsError mocks the Goa endpoint to return an error.
		ReturnsError(error) *Endpoint
	}

	// Expecter adds expectations to a Goa endpoint.
	Expecter interface {
		// WillGet sets expectations on the payload received by the Goa endpoint.
		WillGet(interface{}) *Endpoint
		// WillReturn sets expectations on the result returned by the Goa endpoint.
		WillReturn(interface{}) *Endpoint
		// WillReturnError sets expectations that the Goa endpoint will return an
		// error.
		WillReturnError(error) *Endpoint
		// WillNeverBeCalled sets expectations that the underlying Goa endpoint
		// will never be invoked.
		WillNeverBeCalled()
	}*/

	// Endpoint encapsulates a Goa endpoint and runs any assertions set on the
	// endpoint.
	Endpoint struct {
		// service is the service name.
		service string
		// endpoint is the endpoint name.
		endpoint string
		// ep is the encapsulated endpoint.
		ep goa.Endpoint
		// beforeMock is run before the underlying endpoint is executed.
		// Mock defined by Mocker.Get gets executed here.
		beforeMock *assertion
		// afterMock is run after the beforeMock and circumvents the underlying
		// endpoint. Mock defined by Mocker.Returns or Mocker.ReturnsError gets
		// executed here.
		afterMock *assertion
		// beforeExpect is run before the underlying endpoint is executed.
		// It take precedence over beforeMock.
		beforeExpect *assertion
		// afterExpect is run after the underlying endpoint is executed.
		// It take precedence over afterMock.
		afterExpect *assertion

		mu *sync.Mutex
	}

	// assertion defines the endpoint assertions.
	assertion struct {
		// payload is the payload that the endpoint gets.
		payload interface{}
		// result is the result that the endpoint returns.
		result interface{}
		// err is the error that the endpoint returns.
		err error
		// neverTrigger if true indicates that the assertion must not be triggered.
		neverTrigger bool
		// triggered if true indicates that the assertion was triggered.
		triggered bool

		mu *sync.Mutex
	}
)

// NewEndpoint returns a test endpoint to make assertions on the given
// endpoint.
func NewEndpoint(svc, ep string, e goa.Endpoint) *Endpoint {
	return &Endpoint{
		service:  svc,
		endpoint: ep,
		ep:       e,
		mu:       new(sync.Mutex),
	}
}

// Endpoint returns a test endpoint which wraps the underlying endpoint.
// The test endpoint ensures that any assertions set on the endpoint are met.
// It resets any assertions made on the endpoint after execution.
func (e *Endpoint) Endpoint() goa.Endpoint {
	return func(ctx context.Context, req interface{}) (res interface{}, err error) {
		// assertion errors
		var aerr error

		e.mu.Lock()
		defer func() {
			if aerr == nil {
				if aerr = e.metAssertions(); aerr != nil {
					err = aerr
				}
			} else {
				err = aerr
			}
			e.reset()
			e.mu.Unlock()
		}()

		// run before expectations
		if a := e.beforeExpect; a != nil {
			a.trigger()
			if !reflect.DeepEqual(req, a.payload) {
				aerr = fmt.Errorf("did not receive expected payload")
				return nil, aerr
			}
		} else if a := e.beforeMock; a != nil {
			req = a.payload
			a.trigger()
		}

		if a := e.afterMock; a != nil && e.afterExpect == nil {
			// circumvent the actual endpoint
			a.trigger()
			return a.result, a.err
		}

		// call the actual endpoint
		res, err = e.ep(ctx, req)
		if a := e.afterExpect; a != nil {
			a.trigger()
			if a.result != nil && !reflect.DeepEqual(res, a.result) {
				aerr = fmt.Errorf("did not return expected result")
				return nil, aerr
			}
			if !reflect.DeepEqual(err, a.err) {
				aerr = fmt.Errorf("did not return expected error")
				return nil, aerr
			}
		}
		return
	}
}

// Gets mocks the payload received by the Goa endpoint.
func (e *Endpoint) Gets(payload interface{}) *Endpoint {
	e.mu.Lock()
	e.beforeMock = &assertion{payload: payload, mu: new(sync.Mutex)}
	e.mu.Unlock()
	return e
}

// Returns mocks the result returned by the Goa endpoint.
func (e *Endpoint) Returns(result interface{}) *Endpoint {
	e.mu.Lock()
	e.afterMock = &assertion{result: result, mu: new(sync.Mutex)}
	e.mu.Unlock()
	return e
}

// ReturnsError mocks the Goa endpoint to return an error.
func (e *Endpoint) ReturnsError(err error) *Endpoint {
	e.mu.Lock()
	e.afterMock = &assertion{err: err, mu: new(sync.Mutex)}
	e.mu.Unlock()
	return e
}

// WillGet sets expectations on the payload received by the Goa endpoint.
func (e *Endpoint) WillGet(payload interface{}) *Endpoint {
	e.mu.Lock()
	e.beforeExpect = &assertion{payload: payload, mu: new(sync.Mutex)}
	e.mu.Unlock()
	return e
}

// WillReturn sets expectations on the result returned by the Goa endpoint.
func (e *Endpoint) WillReturn(result interface{}) *Endpoint {
	e.mu.Lock()
	e.afterExpect = &assertion{result: result, mu: new(sync.Mutex)}
	e.mu.Unlock()
	return e
}

// WillReturnError sets expectations that the Goa endpoint will return an
// error.
func (e *Endpoint) WillReturnError(err error) *Endpoint {
	e.mu.Lock()
	e.afterExpect = &assertion{err: err, mu: new(sync.Mutex)}
	e.mu.Unlock()
	return e
}

// WillNeverBeCalled sets expectations that the underlying Goa endpoint will
// never be invoked.
func (e *Endpoint) WillNeverBeCalled() *Endpoint {
	e.mu.Lock()
	e.afterExpect = &assertion{neverTrigger: true, mu: new(sync.Mutex)}
	e.mu.Unlock()
	return e
}

// metAssertions ensures all assertions set on the endpoint are met. It returns
// an error if there are any unmet assertions. Locking is performed by the
// caller.
func (e *Endpoint) metAssertions() error {
	errs := []error{}

	assert := func(a *assertion) {
		if a == nil {
			return
		}
		if err := a.met(); err != nil {
			errs = append(errs, err)
		}
	}

	assert(e.beforeMock)
	assert(e.afterMock)
	assert(e.beforeExpect)
	assert(e.afterExpect)
	if len(errs) > 0 {
		var e string
		for _, er := range errs {
			e += er.Error() + "\n"
		}
		return fmt.Errorf("expectations not met: %s", e)
	}
	return nil
}

// reset resets the mocks and expectations. Locking is performed by the caller.
func (e *Endpoint) reset() {
	e.beforeMock = nil
	e.afterMock = nil
	e.beforeExpect = nil
	e.afterExpect = nil
}

// trigger triggers an assertion.
func (a *assertion) trigger() {
	a.mu.Lock()
	a.triggered = true
	a.mu.Unlock()
}

// met ensures the assertion is met. The assertion is not met if at least one
// of payload, result, or error is not nil.
func (a *assertion) met() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.neverTrigger && a.triggered {
		return fmt.Errorf("expected assertion to not be triggered")
	}
	if !a.triggered {
		return fmt.Errorf("assertion not triggered")
	}
	return nil
}
