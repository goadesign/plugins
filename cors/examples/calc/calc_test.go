package calc

import (
	"log"
	"testing"

	goahttp "goa.design/goa/v3/http"
	"goa.design/plugins/v3/cors/examples/calc/gen/calc"
	calcserver "goa.design/plugins/v3/cors/examples/calc/gen/http/calc/server"
)

func TestMountIndexHTML(t *testing.T) {
	mux := goahttp.NewMuxer()
	e := calc.NewEndpoints(NewCalc(log.Default()))
	s := calcserver.New(e, mux, nil, nil, nil, nil, nil)
	calcserver.Mount(mux, s)
}
