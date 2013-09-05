package gojrest

import (
	"net/http"
	"strings"
)

type RouteWithArg struct {
	Parameter string
	Routes    map[string]HTTPFunc
	Steps     []Stepper
}

func NewRouteWithArg(parameter string, routes map[string]HTTPFunc, steps ...Stepper) *RouteWithArg {
	return &RouteWithArg{
		Parameter: parameter,
		Routes:    routes,
		Steps:     steps,
	}
}

func (r *RouteWithArg) Step(req *http.Request, url string, parameters map[string]string) (bool, interface{}) {
	trimmed_url_parts := strings.SplitN(url[1:], "/", 2)
	if len(trimmed_url_parts) == 0 {
		return false, nil
	}

	value := trimmed_url_parts[0]
	parameters[r.Parameter] = value
	//defer delete(parameters, r.Parameter)

	if len(trimmed_url_parts) == 2 && len(trimmed_url_parts[1]) != 0 {
		// step time
		return UtilStepThroughSteps(req, "/"+trimmed_url_parts[1], parameters, r.Steps)
	}

	if method, exists := r.Routes[req.Method]; exists {
		r := method(req, parameters) // cannot be in the same call as return since defer statements run before the arguments to return are evaluated
		return true, r
	}
	return false, nil
}
