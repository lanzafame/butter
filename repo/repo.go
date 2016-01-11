package repo

import (
	"fmt"
	"time"
	"io"

	"golang.org/x/crypto/ssh"

	"github.com/nanopack/butter/config"
)

type (
	File struct {
		Name string
		Size int64
		IsDir bool
		ModTime time.Time
	}

	Commit struct {
		ID string
		Message string
		AuthorName string
		AuthorEmail string
		Timestamp time.Time
	}

	Command interface {
		Match(string) bool
		Run(string, ssh.Channel) (uint64, error)
	}

	Manager interface {
		Initialize() error
		Commands() []Command
		ListBranches() ([]string, error)
		GetBranch(id string) (string, error)
		ListCommits(branch string, page int) ([]Commit, error)
		GetCommit(id string) (Commit, error)
		ListFiles(commit string) ([]File, error)
		GetFile(commit, id string) (File, error)
		GetFileReader(commit, id string) (io.ReadCloser, error)
	}
)

var (
	availableManagers = map[string]Manager{}
	defaultManager Manager

)

func Register(name string, m Manager) {
	availableManagers[name] = m
}

func Setup() error {
	manager, ok := availableManagers[config.RepoType]
	if !ok {
		return fmt.Errorf("no repo manager found for %s", config.RepoType)
	}
	defaultManager = manager
	return manager.Initialize()
	// setup the default client and call initialize on it once its setup
}

func Commands() []Command {
	return defaultManager.Commands()
}

func ListBranches() ([]string, error) {
	return defaultManager.ListBranches()
}

func GetBranch(id string) (string, error) {
	return defaultManager.GetBranch(id)
}

func ListCommits(branch string, page int) ([]Commit, error) {
	return defaultManager.ListCommits(branch, page)
}

func GetCommit(id string) (Commit, error) {
	return defaultManager.GetCommit(id)
}

func ListFiles(commit string) ([]File, error) {
	return defaultManager.ListFiles(commit)
}

func GetFile(commit, id string) (File, error) {
	return defaultManager.GetFile(commit, id)
}

func GetFileReader(commit, id string) (io.ReadCloser, error) {
	return defaultManager.GetFileReader(commit, id)
}
