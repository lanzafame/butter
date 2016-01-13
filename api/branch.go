package api

import (
	"github.com/nanopack/butter/repo"
	"net/http"
)

func showBranches(rw http.ResponseWriter, req *http.Request) {
	branches, err := repo.ListBranches()
	if err != nil {
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeBody(branches, rw, http.StatusOK)
}

// there arent branch details yet... as far as i know
func showBranchDetails(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte(req.URL.Query().Get(":branch")))
}
