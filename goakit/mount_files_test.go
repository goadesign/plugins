package goakit

import (
	"testing"

	"goa.design/goa/expr"
	httpcodegen "goa.design/goa/http/codegen"
	"goa.design/plugins/goakit/testdata"
)

func TestMountFiles(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string][]string
	}{
		"multi-endpoints": {
			DSL: testdata.MultiEndpointDSL,
			Code: map[string][]string{
				"goakit-mount-handler":     []string{testdata.Endpoint1GoakitMountCode, testdata.Endpoint2GoakitMountCode},
				"goakit-mount-file-server": []string{},
			},
		},
		"file-servers": {
			DSL: testdata.FileServerDSL,
			Code: map[string][]string{
				"goakit-mount-handler":     []string{},
				"goakit-mount-file-server": []string{testdata.File1GoakitMountCode, testdata.File2GoakitMountCode},
			},
		},
		"mixed": {
			DSL: testdata.MixedDSL,
			Code: map[string][]string{
				"goakit-mount-handler":     []string{testdata.MixedMethodGoakitMountCode},
				"goakit-mount-file-server": []string{testdata.MixedFileGoakitMountCode},
			},
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			fs := MountFiles(expr.Root)
			if len(fs) != 1 {
				t.Fatalf("got %d files, expected 1", len(fs))
			}
			for sec, secCode := range c.Code {
				testCode(t, fs[0], sec, secCode)
			}
		})
	}
}
