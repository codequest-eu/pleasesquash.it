package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
)

var (
	binding = flag.String("binding", ":443", "interface and port to bind to")
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
	stateStore := newSecureCookieStore(
		[]byte(os.Getenv("COOKIE_HASH_KEY")),
		[]byte(os.Getenv("COOKIE_BLOCK_KEY")),
	)
	credStore, err := newBoltCredentialsStore(os.Getenv("CRED_STORE"))
	if err != nil {
		glog.Fatal(err)
	}
	oauthManager := newOauthFactory(
		os.Getenv("GITHUB_CLIENT_KEY"),
		os.Getenv("GITHUB_CLIENT_SECRET"),
	)

	h := &handler{
		state: stateStore,
		creds: credStore,
		oauth: oauthManager,
	}
	router := mux.NewRouter()
	router.HandleFunc("/", catchError(serveIndex)).Methods("GET")
	router.HandleFunc("/submit", catchError(h.submit)).Methods("POST")
	router.HandleFunc("/callback", catchError(h.callback)).Methods("GET")
	router.HandleFunc("/webhook", catchError(h.webhook)).Methods("POST")
	glog.Infof("Listening on %s", *binding)
	// TODO: implement TLS as an option.
	glog.Fatal(
		http.ListenAndServeTLS(
			*binding,
			os.Getenv("SSH_CERT"),
			os.Getenv("SSH_KEY"),
			router,
		),
	)
}
