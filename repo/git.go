package repo

import (
	"fmt"
	"io/ioutil"
	"strings"
	"os/exec" 
	"os"
	"io"

	"golang.org/x/crypto/ssh"

	"github.com/gogits/git"

	"github.com/nanopack/butter/config"
)

type (
	gitRepo struct {} // Manager interface
	Push struct{} // Command interface
	Pull struct{} // Command interface
)

func init() {
	Register("git", gitRepo{})
}

func (g gitRepo) Initialize() error {
	os.Mkdir(config.RepoLocation, 0777)
	if _, err := os.Stat(config.RepoLocation + "/info"); os.IsNotExist(err) {
		cmd := exec.Command("git", "init", "--bare")
		cmd.Dir = config.RepoLocation
		return cmd.Run()
	}
	return nil
}

func (g gitRepo) Commands() []Command {
	return []Command{Push{}, Pull{}}
}

func (g gitRepo) ListBranches() ([]string, error) {
	repo, err := git.OpenRepository(config.RepoLocation)
	if err != nil {
		return nil, err
	}
	return repo.GetBranches()
}

func (g gitRepo) GetBranch(id string) (string, error) {
	return id, nil
}

func (g gitRepo) ListCommits(branch string, page int) ([]Commit, error) {
	repo, err := git.OpenRepository(config.RepoLocation)
	if err != nil {
		return nil, err
	}
	commit, err := repo.GetCommitOfBranch(branch)
	if err != nil {
		return nil, err
	}
	list, err := commit.CommitsByRange(page)
	if err != nil {
		return nil, err
	}
	commits := []Commit{}
	for elem := list.Front(); elem != nil; elem = elem.Next() {
		c, ok := elem.Value.(*git.Commit)
		if !ok {
			return nil, fmt.Errorf("the element value is of type %#v", elem.Value)
		}
		com := Commit{
			ID: c.Id.String(), 
			Message: c.CommitMessage,
			AuthorName: c.Author.Name,
			AuthorEmail: c.Author.Email,
			Timestamp: c.Author.When,
		}
		commits = append(commits, com)
	}

	return nil, nil
}

func (g gitRepo) GetCommit(id string) (Commit, error) {
	repo, err := git.OpenRepository(config.RepoLocation)
	if err != nil {
		return Commit{}, err
	}
	c, err := repo.GetCommit(id)
	if err != nil {
		return Commit{}, err
	}
	com := Commit{
			ID: c.Id.String(), 
			Message: c.CommitMessage,
			AuthorName: c.Author.Name,
			AuthorEmail: c.Author.Email,
			Timestamp: c.Author.When,
		}
	return com, nil	
}

func (g gitRepo) ListFiles(commit string) ([]File, error) {
	repo, err := git.OpenRepository(config.RepoLocation)
	if err != nil {
		return nil, err
	}
	tree, err := repo.GetTree(commit) 
	if err != nil {
		return nil, err
	}
	files := []File{}
	for _, entry := range tree.ListEntries() {
		file := File{
			Name: entry.Name(),
			Size: entry.Size(),
			IsDir: entry.IsDir(),
			ModTime: entry.ModTime(),
		}
		files = append(files, file)
	}

	return files, nil
}

func (g gitRepo) GetFile(commit, path string) (File, error) {
	repo, err := git.OpenRepository(config.RepoLocation)
	if err != nil {
		return File{}, err
	}
	tree, err := repo.GetTree(commit) 
	if err != nil {
		return File{}, err
	}
	entry, err := tree.GetTreeEntryByPath(path)
	if err != nil {
		return File{}, err
	}
	file := File{
		Name: entry.Name(),
		Size: entry.Size(),
		IsDir: entry.IsDir(),
		ModTime: entry.ModTime(),
	}

	return file, nil
}

func (g gitRepo) GetFileReader(commit, path string) (io.ReadCloser, error) {
	repo, err := git.OpenRepository(config.RepoLocation)
	if err != nil {
		return nil, err
	}
	tree, err := repo.GetTree(commit) 
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
	originalCommit := getCommit("master")
	code, err := gitShell(ch, ch.Stderr(), command)
	if err == nil {
		newCommit := getCommit("master")
		if newCommit == originalCommit {
			// nothing happened
			return code, nil
		}

		// stream := ch.Stderr()
		// err := nanobox.Deploy(stream, newCommit)

		// switch {
		// case err == templates.ApiUnavailableError:

		// 	// write the original commit back to the master file, so that we
		// 	// can trigger a deploy again without needing new code to be
		// 	// pushed
		// 	name := headName("master")
		// 	ioutil.WriteFile(name, []byte(originalCommit+"\n"), 0600)
		// 	fallthrough

		// case err != nil:

		// 	// we return nil because we have already sent the error across
		// 	// in the nanobox.Deploy function
		// 	return 1, nil
		// }
	}
	return code, err
}

func (pull Pull) Match(command string) bool {
	return strings.HasPrefix(command, "git-send-pack ")
}

func (pull Pull) Run(command string, ch ssh.Channel) (uint64, error) {
	return gitShell(ch, ch.Stderr(), command)
}
func headName(name string) string {
	return config.RepoLocation + "/refs/heads/" + name
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

	err := cmd.Run()
	fmt.Println(err) // logging

	if err != nil {
		// should return the actual exit code
		return 1, err
	}

	return 0, nil
}