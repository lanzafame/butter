package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/nanopack/butter/repo"
)

func listFiles(rw http.ResponseWriter, req *http.Request) {
	comm := req.FormValue("commit")
	fmt.Println("commmmm",comm)
	files, err := repo.ListFiles(comm)
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
