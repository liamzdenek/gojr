package gojr

import "strings"

type APICall struct {
	Method string
	Function HTTPFunc
}

func NewAPICall(method string, function HTTPFunc) *APICall {
	return &APICall{Method: strings.ToUpper(method), Function: function}
}

func (c *APICall) Step(req *Request, url string) (bool, interface{}) {
	if (len(url) == 0 || url == "/") && c.Method == strings.ToUpper(req.Method) /* c.Method is already guaranteed to be upper case */ {
		return true, c.Function(req)
	}
	return false, nil
}
