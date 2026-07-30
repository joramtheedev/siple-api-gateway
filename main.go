package main

import (
	"net/http"
)

func main() {
	//Initialize the api Gateway
	apigateway := openapi.NewAPIGateway()

	//Define routes
	routes := []Route{
		{"POST", "/users", "/services/user-service",
			func(w http.Response, r *http.Reques) {
				apiGateway.proxyUserService(r, w)
			},
		},
	}
	apiGateway.SetRoutes(routes)

	//start the server
	http.ListenAndServe(":8080", apiGateway.Server())
}
