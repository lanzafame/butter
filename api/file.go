package api

import (
	"net/http"
	"io"
	
	"github.com/nanopack/butter/repo"
)

func listFiles(rw http.ResponseWriter, req *http.Request) {
	files, err := repo.ListFiles(req.FormValue("commit"))
	if err != nil {
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	writeBody(files, rw, http.StatusOK)
}

func getFileContents(rw http.ResponseWriter, req *http.Request) {
	reader, err := repo.GetFileReader(req.FormValue("commit"), req.URL.Query().Get(":file"))
	if err != nil {
		rw.Write([]byte(err.Error()))
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(rw, reader)
	return
}
