package gojr

type APICall struct {
	Method string
	Function HTTPFunc
}

func NewAPICall(method string, function HTTPFunc) *APICall {
	return &APICall{Method: method, Function: function}
}

func (c *APICall) Step(req *Request, url string) (bool, interface{}) {
	if (len(url) == 0 || url == "/") && c.Method == req.Method {
		return true, c.Function(req)
	}
	return false, nil
}
