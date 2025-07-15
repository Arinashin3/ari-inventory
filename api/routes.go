package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/register", HandlerRegister)
	mux.HandleFunc("/api/query", HandlerQuery)
	mux.HandleFunc("/api/clear", HandlerClear)
}
