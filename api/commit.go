package api

import (
	"net/http"
	"strconv"

	"github.com/nanopack/butter/repo"
)

func showCommits(rw http.ResponseWriter, req *http.Request) {
	page, _ := strconv.Atoi(req.FormValue("page"))
	commits, err := repo.ListCommits(req.FormValue("page"), page)
	if err != nil {
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeBody(commits, rw, http.StatusOK)
}

func showCommitDetails(rw http.ResponseWriter, req *http.Request) {
	commit, err := repo.GetCommit(req.URL.Query().Get(":commit"))
	if err != nil {
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeBody(commit, rw, http.StatusOK)
}
