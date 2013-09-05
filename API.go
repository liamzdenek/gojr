package gojrest;

import (
	"net/http";
	"encoding/json"
);

type HTTPFunc func(*http.Request, map[string]string) interface{}

type Stepper interface {
	Step(*http.Request, string, map[string]string) (bool, interface{})
}

type API struct {
	Steps []Stepper
}

func NewAPI(steps ...Stepper) *API {
	return &API{
		Steps: steps,
	}
}

func (api *API) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	didroute, response := UtilStepThroughSteps(req, "/"+req.URL.Path, map[string]string{}, api.Steps)
	if didroute {
		j, _ := json.Marshal(response)
		res.Write(j)
		return
	}
	res.WriteHeader(404)
	j, _ := json.Marshal(struct{}{})
	res.Write(j)
}

func UtilStepThroughSteps(req *http.Request, url string, parameters map[string]string, steps []Stepper) (bool, interface{}) {
	for _, step := range steps {
		didroute, response := step.Step(req, url, parameters)
		if didroute {
			return true, response // did match
		}
	}
	return false, nil // had no match
}
