package handler

import "net/http"

func proxyUserService(r *http.Request, w http.ResponseWriter) {
	//forward the request to the initial service
	internalServiceURL := r.FormValue("internal-service-url")
}
