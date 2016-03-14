package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

var (
	binding = flag.String("binding", ":1983", "interface and port to bind to")
)

type fallibleHandler func(w http.ResponseWriter, r *http.Request) error

func catchError(fn fallibleHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			glog.Warningf("Error for %q: %v", r.URL.String(), err)
			http.Error(w, "You broke the internet 🐮💩😱", http.StatusInternalServerError)
		}
	}
}

func main() {
	flag.Parse()
	h := &handler{} // TODO
	router := mux.NewRouter()
	router.HandleFunc("/", catchError(serveIndex)).Methods("GET")
	router.HandleFunc("/submit", catchError(h.submit)).Methods("POST")
	router.HandleFunc("/callback", catchError(h.callback)).Methods("GET")
	router.HandleFunc("/webhook", catchError(h.webhook)).Methods("POST")
	glog.Infof("Listening on %s", *binding)
	// TODO: implement TLS as an option.
	glog.Fatal(http.ListenAndServe(*binding, router))
}
