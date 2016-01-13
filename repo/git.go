package repo

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/shazow/go-git"

	"github.com/nanopack/butter/config"
	"github.com/nanopack/butter/deploy"
)

type (
	gitRepo struct{} // Manager interface
	Push    struct{} // Command interface
	Pull    struct{} // Command interface
)

func init() {
	Register("git", gitRepo{})
}

func (g gitRepo) repo() (*git.Repository, error) {
	return git.OpenRepository(config.RepoLocation + "/live.git")
}

func (g gitRepo) Initialize() error {
	live := config.RepoLocation + "/live.git"
	os.MkdirAll(live, 0777)
	if _, err := os.Stat(live + "/info"); os.IsNotExist(err) {
		cmd := exec.Command("git", "init", "--bare")
		cmd.Dir = live
		return cmd.Run()
	}
	return nil
}

func (g gitRepo) Commands() []Command {
	return []Command{Push{}, Pull{}}
}

func (g gitRepo) ListBranches() ([]string, error) {
	repo, err := g.repo()
	if err != nil {
		return nil, err
	}
	return repo.GetBranches()
}

func (g gitRepo) GetBranch(id string) (string, error) {
	return id, nil
}

func (g gitRepo) ListCommits(branch string, page int) ([]Commit, error) {
	repo, err := g.repo()
	if err != nil {
		config.Log.Debug("getting repo")
		return nil, err
	}
	commit, err := repo.GetCommitOfBranch(branch)
	if err != nil {
		config.Log.Debug("getting commit")
		return nil, err
	}
	list, err := repo.CommitsBefore(commit.Id.String())
	if err != nil {
		config.Log.Debug("getting commits")
		return nil, err
	}
	commits := []Commit{}
	elem := list.Front()
	paging := true
	if page <= 0 {
		paging = false
	}
	page = page - 1
	for i := 0; i < list.Len() && elem != nil; i++ {
		if !paging || i / 100 == page {
			c, ok := elem.Value.(*git.Commit)
			if !ok {
				return nil, fmt.Errorf("the element value is of type %#v", elem.Value)
			}
			com := Commit{
				ID:          c.Id.String(),
				Message:     c.CommitMessage,
				AuthorName:  c.Author.Name,
				AuthorEmail: c.Author.Email,
				Timestamp:   c.Author.When,
			}
			commits = append(commits, com)

		}
		elem = elem.Next()	
	}
	for ; elem != nil;  {
	}
	return commits, nil
}

func (g gitRepo) GetCommit(id string) (Commit, error) {
	repo, err := g.repo()
	if err != nil {
		return Commit{}, err
	}
	c, err := repo.GetCommit(id)
	if err != nil {
		return Commit{}, err
	}
	com := Commit{
		ID:          c.Id.String(),
		Message:     c.CommitMessage,
		AuthorName:  c.Author.Name,
		AuthorEmail: c.Author.Email,
		Timestamp:   c.Author.When,
	}
	return com, nil
}

func (g gitRepo) ListFiles(commit string) ([]File, error) {
	repo, err := g.repo()
	if err != nil {
		return nil, err
	}
	fmt.Println("repo", repo)
	c, err := repo.GetCommit(commit)
	if err != nil {
		return nil, err
	}

	tree, err := repo.GetTree(c.TreeId().String())
	if err != nil {
		return nil, err
	}
	fmt.Println("tree", tree)
	files := []File{}
	walker := func(path string, te *git.TreeEntry, err error) error {
		if err != nil {
			return err
		}
		file := File{
			Name:    path,
			Size:    te.Size(),
			IsDir:   te.IsDir(),
			ModTime: te.ModTime(),
		}
		files = append(files, file)
		return nil
	}
	tree.Walk(walker)

	// scanner, err := tree.Scanner()
	// if err != nil {
	// 	return nil, err
	// }

	// files := []File{}
	// for scanner.Scan() {
	// 	fmt.Println(scanner.TreeEntry())
	// 	file := File{
	// 		Name:    scanner.TreeEntry().Name(),
	// 		Size:    scanner.TreeEntry().Size(),
	// 		IsDir:   scanner.TreeEntry().IsDir(),
	// 		ModTime: scanner.TreeEntry().ModTime(),
	// 	}
	// 	files = append(files, file)
	// }	

	return files, nil
}

func (g gitRepo) GetFile(commit, path string) (File, error) {
	repo, err := g.repo()
	if err != nil {
		return File{}, err
	}
	c, err := repo.GetCommit(commit)
	if err != nil {
		return File{}, err
	}

	tree, err := repo.GetTree(c.TreeId().String())
	if err != nil {
		return File{}, err
	}
	entry, err := tree.GetTreeEntryByPath(path)
	if err != nil {
		return File{}, err
	}
	file := File{
		Name:    entry.Name(),
		Size:    entry.Size(),
		IsDir:   entry.IsDir(),
		ModTime: entry.ModTime(),
	}

	return file, nil
}

func (g gitRepo) GetFileReader(commit, path string) (io.ReadCloser, error) {
	fmt.Println("thisisafadsfasfsdaf",path)
	repo, err := g.repo()
	if err != nil {
		fmt.Println("repo")
		return nil, err
	}
	c, err := repo.GetCommit(commit)
	if err != nil {
		return nil, err
	}

	tree, err := repo.GetTree(c.TreeId().String())
	if err != nil {
		return nil, err
	}
	blob, err := tree.GetBlobByPath(path)
	if err != nil {
		return nil, err
	}

	return blob.Data()
}

func (push Push) Match(command string) bool {
	return strings.HasPrefix(command, "git-receive-pack ")
}

func (push Push) Run(command string, ch ssh.Channel) (uint64, error) {
	//TODO make "master" be dynamic
	code, err := gitShell(ch, ch.Stderr(), command)
	if err == nil {
		newCommit := getCommit("master")
		stream := ch.Stderr()
		err = deploy.Run(stream, newCommit)
		if err != nil {
			return 1, err
		}
	}
	return code, err
}

func (pull Pull) Match(command string) bool {
	return strings.HasPrefix(command, "git-send-pack") || strings.HasPrefix(command, "git-upload-pack")
}

func (pull Pull) Run(command string, ch ssh.Channel) (uint64, error) {
	return gitShell(ch, ch.Stderr(), command)
}
func headName(name string) string {
	return config.RepoLocation + "/live.git/refs/heads/" + name
}
func getCommit(name string) string {
	file := headName(name)
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(bytes), "\n\r")
}

func gitShell(duplex io.ReadWriter, errStream io.Writer, command string) (uint64, error) {
	cmd := exec.Command("git", "shell", "-c", command)
	cmd.Dir = config.RepoLocation

	cmd.Stdout = duplex
	cmd.Stderr = errStream
	cmd.Stdin = duplex
	config.Log.Debug("running command %+v", cmd)

	err := cmd.Run()
	fmt.Println(err) // logging

	if err != nil {
		// should return the actual exit code
		return 1, err
	}

	return 0, nil
}
