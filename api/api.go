package api

import (
	"net/http"
	
	"github.com/gorilla/pat"
	"github.com/nanobox-io/nanoauth"
	"github.com/nanopack/butter/config"
	"encoding/json"
)

func Start() error {
	router := pat.New()

	router.Get("/branches", showBranches)
	router.Get("/branches/{branch}", showBranchDetails)
	router.Get("/commits", showCommits)
	router.Get("/commits/{commit}", showCommitDetails)
	router.Get("/files", listFiles)
	router.Get("/files/{file}", getFileContents)

	// blocking...
	return nanoauth.ListenAndServeTLS(config.HttpListenAddress, config.Token, router)
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