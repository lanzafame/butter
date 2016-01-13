package api

import (
	"net/http"

	"encoding/json"
	"github.com/gorilla/pat"
	"github.com/nanobox-io/nanoauth"
	"github.com/nanopack/butter/config"
)

func Start() error {
	router := pat.New()

	router.Get("/branches/{branch}", handleRequest(showBranchDetails))
	router.Get("/branches", handleRequest(showBranches))
	router.Get("/commits/{commit}", handleRequest(showCommitDetails))
	router.Get("/commits", handleRequest(showCommits))
	router.Get("/files/{file:.*}", handleRequest(getFileContents))
	router.Get("/files", handleRequest(listFiles))

	// blocking...
	config.Log.Info("Api Listening on %s", config.HttpListenAddress)
	return nanoauth.ListenAndServeTLS(config.HttpListenAddress, config.Token, router)
}

// handleRequest
func handleRequest(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {

		config.Log.Debug(`
Request:
--------------------------------------------------------------------------------
%+v

`, req)

		//
		fn(rw, req)

		config.Log.Debug(`
Response:
--------------------------------------------------------------------------------
%+v

`, rw)
	}
}

// writeBody
func writeBody(v interface{}, rw http.ResponseWriter, status int) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	rw.Write(b)

	return nil
}
