package api

import (
	"net/http"
	"strings"
)

func ListenAndServe(addr string, server Server, redirector Redirector) error {
	http.HandleFunc("/", mainHandler(server, redirector))
	return http.ListenAndServe(addr, nil)
}

func mainHandler(server Server, redirector Redirector) http.HandlerFunc {
	routes := map[string]http.HandlerFunc{
		"/status": StatusHandler,
		"/scan":   API(server.Scan),
		"/get":    API(server.Get),
		"/put":    API(server.Put),
		"/delete": API(server.Delete),
	}
	redirect := RedirectHandler(redirector, "/c/")
	return func(w http.ResponseWriter, r *http.Request) {
		if h, ok := routes[r.URL.Path]; ok {
			h(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/c") {
			redirect(w, r)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}
}
