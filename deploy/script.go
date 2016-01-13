package deploy

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/nanopack/butter/config"
)

type (
	script struct{}
)

func init() {
	Register("script", script{})
}

func (s script) Initialize() error {
	file, err := os.Open(config.DeployLocation)
	if err != nil {
		return fmt.Errorf("deploy init: %+v", err)
	}
	fi, err := file.Stat()
	if err != nil {
		return fmt.Errorf("deploy init: %+v", err)
	}
	if fi.IsDir() {
		return fmt.Errorf("file given is a directory")
	}
	return nil
}

func (s script) Run(stream io.Writer, id string) error {
	cmd := exec.Command(config.DeployLocation, id)
	cmd.Stdout = stream
	cmd.Stderr = stream

	return cmd.Run()
}
