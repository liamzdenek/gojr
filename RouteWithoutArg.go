package gojrest

import (
	"net/http"
	"strings"
)

type RouteWithoutArg struct {
	Prefix string
	Routes map[string]HTTPFunc
	Steps  []Stepper
}

func NewRouteWithoutArg(prefix string, routes map[string]HTTPFunc, steps ...Stepper) *RouteWithoutArg {
	return &RouteWithoutArg{
		Prefix: "/" + prefix,
		Routes: routes,
		Steps:  steps,
	}
}

func (r *RouteWithoutArg) Step(req *http.Request, url string, parameters map[string]string) (bool, interface{}) {
	trimmed_url := strings.TrimPrefix(url, r.Prefix)
	if len(trimmed_url) > len(url) {
		return false, nil // didn't match this at all
	}

	if len(trimmed_url) == 0 || (len(trimmed_url) == 1 && trimmed_url == "/") {
		if method, exists := r.Routes[req.Method]; exists {
			return true, method(req, parameters)
		}
		return false, nil // post to us directly
	}
	return UtilStepThroughSteps(req, trimmed_url, parameters, r.Steps)
}
