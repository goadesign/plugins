package goakit

import (
	"testing"

	"goa.design/goa/v3/expr"
	httpcodegen "goa.design/goa/v3/http/codegen"
	"goa.design/plugins/v3/goakit/testdata"
)

func TestMountFiles(t *testing.T) {
	cases := map[string]struct {
		DSL  func()
		Code map[string][]string
		Path string
	}{
		"multi-endpoints": {
			DSL: testdata.MultiEndpointDSL,
			Code: map[string][]string{
				"goakit-mount-handler":     {testdata.Endpoint1GoakitMountCode, testdata.Endpoint2GoakitMountCode},
				"goakit-mount-file-server": {},
			},
			Path: "gen/http/multi_endpoint_service/kitserver/mount.go",
		},
		"file-servers": {
			DSL: testdata.FileServerDSL,
			Code: map[string][]string{
				"goakit-mount-handler":     {},
				"goakit-mount-file-server": {testdata.File1GoakitMountCode, testdata.File2GoakitMountCode},
			},
			Path: "gen/http/file_server_service/kitserver/mount.go",
		},
		"mixed": {
			DSL: testdata.MixedDSL,
			Code: map[string][]string{
				"goakit-mount-handler":     {testdata.MixedMethodGoakitMountCode},
				"goakit-mount-file-server": {testdata.MixedFileGoakitMountCode},
			},
			Path: "gen/http/mixed_service/kitserver/mount.go",
		},
		"goifyable": {
			DSL: testdata.GoifyableServiceDSL,
			Code: map[string][]string{
				"goakit-mount-handler":     {testdata.GoifyableMethodGoakitMountCode},
				"goakit-mount-file-server": {},
			},
			Path: "gen/http/goifyable_service/kitserver/mount.go",
		},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			httpcodegen.RunHTTPDSL(t, c.DSL)
			fs := MountFiles(expr.Root)
			if len(fs) != 1 {
				t.Fatalf("got %d files, expected 1", len(fs))
			}
			f := fs[0]
			if f.Path != c.Path {
				t.Errorf("got path %q, expected %q", f.Path, c.Path)
			}
			for sec, secCode := range c.Code {
				testCode(t, f, sec, secCode)
			}
		})
	}
}
