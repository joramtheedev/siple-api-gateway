package models

import "net/http"

type Route struct {
	Method  string
	Path    string
	Service string
	Handler func(http.ResponseWriter, *http.Request)
}
