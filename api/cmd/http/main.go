package main

import (
	"net/http"

	handler_http "github.com/IndominusByte/farmacare-be/api/cmd/http/handler"
)

func main() {
	// mount router
	r := handler_http.CreateNewServer()
	if err := r.MountHandlers(); err != nil {
		panic(err)
	}

	http.ListenAndServe(":3000", r.Router)
}
