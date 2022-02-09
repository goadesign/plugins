// Supporting elements shared among tests in this package.

package goakit

import (
	"testing"

	"goa.design/goa/v3/codegen"
)

func requireImport(t *testing.T, f *codegen.File, importPath string) {
	t.Helper()
	var found bool
	hd := f.Section("source-header")[0].Data.(map[string]interface{})
	for _, i := range hd["Imports"].([]*codegen.ImportSpec) {
		if i.Path == importPath {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("%q does not import %q", f.Path, importPath)
	}
}
