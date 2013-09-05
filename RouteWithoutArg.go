package gojr

import (
	"strings"
)

type RouteWithoutArg struct {
	Prefix string
	Steps  []Stepper
}

func NewRouteWithoutArg(prefix string, steps ...Stepper) *RouteWithoutArg {
	return &RouteWithoutArg{
		Prefix: "/" + prefix,
		Steps:  steps,
	}
}

func (r *RouteWithoutArg) Step(req *Request, url string) (bool, interface{}) {
	trimmed_url := strings.TrimPrefix(url, r.Prefix)
	if len(trimmed_url) > len(url) {
		return false, nil // didn't match this at all
	}
	return UtilStepThroughSteps(req, trimmed_url, r.Steps)
}
