gojr
====

gojr (pronounced "Go Junior" or simply "Junior") is a microframework for rapid development of JSON APIs in Go. gojr is simply a package that provides a net/http.Handler, and can be used in any existing net/http project with ease.

gojr is under frequent development and is not yet stable, adequately tested, or adequately documented. Usage in production systems is not recommended.

Example
-------

See "examples/main.go"

How gojr Works
--------------

gojr maintains a tree where each node is a portion of the requested URL's path. There are four main classes used when defining a gojr program
* API
** API contains a Tree, and is exclusively a root node.
* RouteWithoutArgument
** RouteWithoutArgument contains a static segment of a URL's path.
** eg, "/user"
* RouteWithArgument
** RouteWithArgument contains a dynamic segment of a URL's path.
** eg. "/:username"
* APICall
** APICall is the endpoint through the tree. This will actually trigger your logic if we have reached the end of the path and the method specified in the APICall matches the method of the request.

A typical gojr program may contain a tree that looks something like this:

```go
// define the gojr structure
api := gojr.NewAPI(
	gojr.NewRouteWithoutArg("user", // parses "/user"
		gojr.NewRouteWithArg("username", parses "/:username" [which could be any string]
			
			// curl -X GET localhost:8080/user/anything/
			gojr.NewAPICall("GET", func(req *gojr.Request) interface{} { // and register and endpoint for the request tree
				return struct{ Username string }{ req.Parameters["username"] }
				// request returns '{"Username":"WhateverYouPutInThePath"}'
			}),
		),
	),
);

// start up the HTTP server
prefix := "api"
http.Handle(prefix, http.StripPrefix(prefix, api))
http.ListenAndServe(":8080")
```
