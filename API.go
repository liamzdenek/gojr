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

type ErrorValueMissing struct {
	Key string
}

func (r *Request) MustFormValue(k string) (v string) {
	v = r.FormValue(k)
	if v == "" {
		panic(ErrorValueMissing{Key: k})
	}
	return
}

func (r *Request) MustPostFormValue(k string) (v string) {
	v = r.PostFormValue(k)
	if v == "" {
		panic(ErrorValueMissing{Key: k})
	}
	return
}

type API struct {
	Steps []Stepper

	// used as the response code when an ErrorValueMissing is thrown
	ErrorCode_ValueMissing  int
	ErrorCode_InternalPanic int
	ErrorCode_RouteNotFound int
}

func NewAPI(steps ...Stepper) *API {
	return &API{
		Steps: steps,
		ErrorCode_ValueMissing:  400,
		ErrorCode_InternalPanic: 500,
		ErrorCode_RouteNotFound: 404,
	}
}

func (api *API) ServeHTTP(res http.ResponseWriter, http_req *http.Request) {
	req := &Request{*http_req, map[string]string{}}

	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case ErrorValueMissing:
				e := r.(ErrorValueMissing)
				res.WriteHeader(api.ErrorCode_ValueMissing)
				j, _ := json.Marshal(struct {
					Error                      bool
					ErrorMessage, ParameterKey string
				}{true, "A required parameter was not provided", e.Key})
				res.Write(j)
			default:
				res.WriteHeader(api.ErrorCode_InternalPanic)
				j, _ := json.Marshal(struct {
					Error        bool
					ErrorMessage string
				}{true, "There was an internal error parsing your request"})
				res.Write(j)
				debug.PrintStack()
			}
		}
	}()

	didroute, response := UtilStepThroughSteps(req, "/"+req.URL.Path, api.Steps)
	if didroute {
		j, _ := json.Marshal(response)
		res.Write(j)
		return
	}
	res.WriteHeader(api.ErrorCode_RouteNotFound)
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
