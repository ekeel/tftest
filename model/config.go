package model

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Config struct {
	Tests []*Test `hcl:"test,block"`
}

func Parse(filepath string) (config Config, err error) {
	config = Config{}

	parser := hclparse.NewParser()

	hclFile, pDiags := parser.ParseHCLFile(filepath)
	if pDiags.HasErrors() {
		return config, fmt.Errorf(pDiags.Error())
	}

	dDiags := gohcl.DecodeBody(hclFile.Body, nil, &config)
	if dDiags.HasErrors() {
		return config, fmt.Errorf(dDiags.Error())
	}

	return config, nil
}
