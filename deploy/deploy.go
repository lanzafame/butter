package deploy

import (
	"github.com/nanopack/butter/config"
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
	defaultDeploy    Deployer
)

func Register(name string, m Deployer) {
	availableDeploys[name] = m
}

func Setup() error {
	deploy, ok := availableDeploys[config.DeployType]
	if ok {
		config.Log.Info("setting up deploy(%s) location: %s", config.DeployType, config.DeployLocation)
		defaultDeploy = deploy
		return deploy.Initialize()
	}
	return nil
}

func Run(stream io.Writer, id string) error {
	config.Log.Debug("running deploy %+v", defaultDeploy)
	if defaultDeploy != nil {
		return defaultDeploy.Run(stream, id)
	}
	return nil
}
