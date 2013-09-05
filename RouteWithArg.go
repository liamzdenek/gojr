package gojr

import (
	"strings"
)

type RouteWithArg struct {
	Parameter string
	Steps     []Stepper
}

func NewRouteWithArg(parameter string, steps ...Stepper) *RouteWithArg {
	return &RouteWithArg{
		Parameter: parameter,
		Steps:     steps,
	}
}

func (r *RouteWithArg) Step(req *Request, url string) (bool, interface{}) {
	trimmed_url_parts := strings.SplitN(url[1:], "/", 2)
	if len(trimmed_url_parts) == 0 {
		return false, nil
	}

	if len(trimmed_url_parts) == 1 {
		trimmed_url_parts = append(trimmed_url_parts, "")
	}

	value := trimmed_url_parts[0]
	req.Parameters[r.Parameter] = value
	defer delete(req.Parameters, r.Parameter)

	found, response := UtilStepThroughSteps(req, "/"+trimmed_url_parts[1], r.Steps) // must not be in the same instruction as the return. defers are ran before the return parameters are evaluated
	return found, response
}
