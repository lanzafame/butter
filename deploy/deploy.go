package deploy

import (
	"io"
)

type (
	Deployer interface {
		Initialize() error
		Run(stream io.Writer, id string) error 
	}
)

var (
	availableDeploys = map[string]Deployer{}
	defaultDeploy Deployer
)

func Register(name string, m Deployer) {
	availableDeploys[name] = m
}

func Setup() error {
	deploy, ok := availableDeploys[config.RepoType]
	if ok {
		defaultDeploy = deploy
		return deploy.Initialize()
	}
	return nil
}

func Run(stream io.Writer, id string) error  {
	return defaultDeploy.Run(stream, id)
}