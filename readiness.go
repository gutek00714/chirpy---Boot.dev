package main

import "net/http"

func healthzHandlerFunction(w http.ResponseWriter, r *http.Request) {
	// set content-type header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	// set the status code
	w.WriteHeader(http.StatusOK)

	// write the body text
	myBytes := []byte("OK")
	w.Write(myBytes)
}
