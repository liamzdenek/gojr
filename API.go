package gojr

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
)

type HTTPFunc func(*Request) interface{}

type Stepper interface {
	Step(*Request, string) (bool, interface{})
}

type Request struct {
	http.Request
	Parameters map[string]string
}

type API struct {
	Steps []Stepper
}

func NewAPI(steps ...Stepper) *API {
	return &API{
		Steps: steps,
	}
}

func (api *API) ServeHTTP(res http.ResponseWriter, http_req *http.Request) {
	req := &Request{*http_req, map[string]string{}}

	defer func() {
		if r := recover(); r != nil {
			res.WriteHeader(500);
			j, _ := json.Marshal(struct{Error bool; ErrorMessage string}{true,"There was an internal error parsing your request"})
			res.Write(j);
			debug.PrintStack()
		}
	}()

	didroute, response := UtilStepThroughSteps(req, "/"+req.URL.Path, api.Steps)
	if didroute {
		j, _ := json.Marshal(response)
		res.Write(j)
		return
	}
	res.WriteHeader(404)
	j, _ := json.Marshal(struct{}{})
	res.Write(j)
}

func UtilStepThroughSteps(req *Request, url string, steps []Stepper) (bool, interface{}) {
	for _, step := range steps {
		didroute, response := step.Step(req, url)
		
		if didroute {
			return true, response // did match
		}
	}
	return false, nil // had no match
}
