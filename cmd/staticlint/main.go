package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"

	"github.com/e-faizov/yibana/cmd/staticlint/checkexit"
)

func main() {
	var checks []*analysis.Analyzer

	checks = append(checks, passes()...)

	for _, v := range staticcheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	checks = append(checks, stylecheck.Analyzers[0].Analyzer)
	checks = append(checks, stylecheck.Analyzers[1].Analyzer)

	checks = append(checks, checkexit.Analyzer)

	multichecker.Main(
		checks...,
	)
}
