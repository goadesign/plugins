package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"goa.design/goa/v3/codegen"
	"goa.design/goa/v3/eval"
	"goa.design/plugins/v3/model"
	"goa.design/plugins/v3/model/expr"

	. "goa.design/plugins/v3/model/dsl"
	testdata "goa.design/plugins/v3/model/testdata/design"
)

func TestModel(t *testing.T) {
	cases := []struct {
		Name     string
		Services func()
		System   string
		Err      string
	}{
		{
			Name:     "exact match",
			Services: testdata.OneService("OneContainer"),
			System:   "OneContainer",
		},
		{
			Name:     "exact match with multiple services",
			Services: testdata.TwoServices("TwoContainers"),
			System:   "TwoContainers",
		},
		{
			Name:     "exact match with excluded tags",
			Services: testdata.OneService("TwoContainers", func() { ModelExcludedTags("A") }),
			System:   "TwoContainers",
		},
		{
			Name:     "exact match with container format",
			Services: testdata.OneService("FormattedContainer", func() { ModelContainerFormat("C %s") }),
			System:   "FormattedContainer",
		},
		{
			Name:     "exact match with explicit container name",
			Services: testdata.OtherService("OneContainer", ModelComplete, func() { ModelContainer("A Service") }),
			System:   "OneContainer",
		},
		{
			Name:     "exact match with model none",
			Services: testdata.ThreeServices("OneContainer", nil, ModelNone),
		},
		{
			Name:     "exact match with default format",
			Services: testdata.OneService("OneContainer", func() { ModelContainerFormat("") }),
			System:   "OneContainer",
		},
		{
			Name:     "No package name",
			Services: testdata.NoPackage(),
			Err:      "model plugin requires a model package path, use the 'Model' function to set it",
		},
		{
			Name:     "No system name",
			Services: testdata.OneService(""),
			Err:      "model plugin requires a system name, use the 'Model' function to set it",
		},
		{
			Name:     "Invalid system name",
			Services: testdata.OneService("Invalid"),
			System:   "OneContainer",
			Err:      "system \"Invalid\" not found",
		},
		{
			Name:     "missing container",
			Services: testdata.TwoServices("OneContainer"),
			System:   "OneContainer",
			Err:      "service B has no corresponding container in the model",
		},
		{
			Name:     "multiple missing containers",
			Services: testdata.ThreeServices("OneContainer"),
			System:   "OneContainer",
			Err:      "services B, C have no corresponding container in the model",
		},
		{
			Name:     "missing service",
			Services: testdata.OneService("TwoContainers", ModelComplete),
			System:   "TwoContainers",
			Err:      "container B Service has no corresponding service in the design",
		},
		{
			Name:     "multiple missing services",
			Services: testdata.OneService("ThreeContainers", ModelComplete),
			System:   "ThreeContainers",
			Err:      "containers B Service, C Service have no corresponding service in the design",
		},
		{
			Name:     "missing container and service",
			Services: testdata.OneService("OtherContainer", ModelComplete),
			System:   "OtherContainer",
			Err:      "service A has no corresponding container in the model, and container C Service has no corresponding service in the design",
		},
		{
			Name:     "multiple missing containers and services",
			Services: testdata.TwoServices("OtherContainers", ModelComplete),
			System:   "OtherContainers",
			Err:      "services A, B have no corresponding container in the model, and containers C Service, D Service have no corresponding service in the design",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			expr.InitRoot()
			root := codegen.RunDSL(t, c.Services)
			_, err := model.Validate("model", []eval.Root{root, expr.Root}, nil)
			if c.Err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), c.Err)
				return
			}
			require.NoError(t, err)
		})
	}
}
