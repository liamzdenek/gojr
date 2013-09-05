package main

import (
	"./gojr"
	"net/http"
)

func main() {
	api := gojr.NewAPI(
		gojr.NewRouteWithoutArg("user",
			// curl -X POST localhost:8080/api/user/
			gojr.NewAPICall("POST", func(req *http.Request, parameters map[string]string) (interface{}) {
				return struct{Hello string; Info string;}{"world", "this would register a user",}
			}),
			gojr.NewRouteWithArg("username",
				// curl localhost:8080/api/user/anystringhere/
				gojr.NewAPICall("GET", func(req *http.Request, parameters map[string]string) (interface{}) {
					return struct{Info string}{"This would return user data for "+parameters["username"]};
				}),
				gojr.NewRouteWithoutArg("login",
					// curl -X POST localhost:8080/api/user/anystringhere/login
					gojr.NewAPICall("POST", func(req *http.Request, parameters map[string]string) (interface{}) {
						return struct{Info string}{"This would log you in to user: "+parameters["username"]}
					}),
				),
			),
		),
	)

	prefix := "/api/"
	http.Handle(prefix, http.StripPrefix(prefix, api))

	http.ListenAndServe(":8080", nil)
}
