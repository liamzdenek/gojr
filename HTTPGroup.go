package gojrest;

func NewHTTPGroup(items ...HTTPItem) map[string]HTTPFunc {
	httpgroup := map[string]HTTPFunc{}
	for _, item := range items {
		httpgroup[item.Method] = item.Function
	}
	return httpgroup
}

type HTTPItem struct {
	Method   string
	Function HTTPFunc
}

func NewHTTPItem(method string, function HTTPFunc) HTTPItem {
	return HTTPItem{
		Method:   method,
		Function: function,
	}
}

