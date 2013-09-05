package gojr

import (
	"net/http"
	"fmt"
)

type APICall struct {
	Method string
	Function HTTPFunc
}

func NewAPICall(method string, function HTTPFunc) *APICall {
	return &APICall{Method: method, Function: function}
}

func (c *APICall) Step(req *http.Request, url string, parameters map[string]string) (bool, interface{}) {
	fmt.Printf("Checking function call - url: %v\n", url);
	if (len(url) == 0 || url == "/") && c.Method == req.Method {
		fmt.Printf("match - %v\n", parameters);
		return true, c.Function(req, parameters)
	}
	return false, nil
}
