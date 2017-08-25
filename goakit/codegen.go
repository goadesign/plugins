package goakit

import (
	"fmt"

	"goa.design/goa/codegen/generator"
)

func init() {
	generator.Generators = generators
}

func generators(cmd string) ([]generator.Genfunc, error) {
	switch cmd {
	case "gen":
		return []generator.Genfunc{
			Service,
			Transport,
			generator.OpenAPI,
		}, nil
	case "example":
		return []generator.Genfunc{
			Example,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported command %q", cmd)
	}
}
