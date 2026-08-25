// This file verifies that the Zap logger plugin changes the starter service
// constructor and every generated call to that constructor together.
package zaplogger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/codegen/generator"
	"goa.design/goa/v3/eval"

	"goa.design/plugins/v3/zaplogger/testdata"
)

func TestExamplePassesZapLoggerToServiceConstructor(t *testing.T) {
	root := codegen.RunDSL(t, testdata.SimpleServiceDSL)
	files, err := generator.Example("generated.local/gen", []eval.Root{root})
	require.NoError(t, err)
	files, err = UpdateExample("generated.local/gen", []eval.Root{root}, files)
	require.NoError(t, err)

	dir := t.TempDir()
	for _, file := range files {
		_, err := file.Render(dir)
		require.NoError(t, err)
	}

	serviceSource, err := os.ReadFile(filepath.Join(dir, "simple_service.go"))
	require.NoError(t, err)
	require.Contains(t, string(serviceSource), "func NewSimpleService(logger *zap.SugaredLogger)")

	var mainSource []byte
	err = filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Name() == "main.go" && strings.Contains(path, "cmd") {
			source, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if strings.Contains(string(source), "Initialize the services") {
				mainSource = source
			}
		}
		return nil
	})
	require.NoError(t, err)
	require.NotEmpty(t, mainSource)
	require.Contains(t, string(mainSource), "simpleServiceSvc = testapi.NewSimpleService(zlog)")
}
