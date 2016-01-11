package config

import (
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/ghodss/yaml"
	"github.com/jcelliott/lumber"
)

var (
	SshListenAddress  string
	HttpListenAddress string
	KeyPath           string
	RepoType					string
	RepoLocation      string
	KeyAuthType				string
	KeyAuthLocation   string
	PassAuthType			string
	PassAuthLocation  string
	DeployType				string
	DeployLocation	  string
	Token             string
	Log            	  lumber.Logger
)

func AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&SshListenAddress, "ssh-address", "", ":2222", "[server] SshListenAddress")
	cmd.Flags().StringVarP(&HttpListenAddress, "http-address", "", ":8080", "[server] HttpListenAddress")
	cmd.Flags().StringVarP(&KeyPath, "key-path", "", "", "[server] KeyPath")
	cmd.Flags().StringVarP(&RepoType, "repo-type", "", "", "[server] RepoType")
	cmd.Flags().StringVarP(&RepoLocation, "repo-location", "", "", "[server] RepoLocation")
	cmd.Flags().StringVarP(&KeyAuthType, "key-auth-type", "", "", "[server] KeyAuthType")
	cmd.Flags().StringVarP(&KeyAuthLocation, "key-auth-location", "", "", "[server] KeyAuthLocation")
	cmd.Flags().StringVarP(&PassAuthType, "pass-auth-type", "", "", "[server] PassAuthType")
	cmd.Flags().StringVarP(&PassAuthLocation, "pass-auth-location", "", "", "[server] PassAuthLocation")
	cmd.Flags().StringVarP(&DeployType, "pass-auth-type", "", "", "[server] DeployType")
	cmd.Flags().StringVarP(&DeployLocation, "pass-auth-location", "", "", "[server] DeployLocation")
	cmd.PersistentFlags().StringVarP(&Token, "token", "", "secret"	, "Token security")
}

func Parse(configFile string) {
	c := map[string]string{}

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		Log = lumber.NewConsoleLogger(lumber.INFO)
		Log.Error("unable to read config file: %v\n", err)
	}
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		Log.Error("err parsing config file: %v\n", err)
		Log.Error("falling back to default values")
	}

	Log = lumber.NewConsoleLogger(lumber.LvlInt(c["log_level"]))
	if c["ssh_listen_address"] != "" {
		SshListenAddress = c["ssh_listen_address"]
	}
	if c["http_listen_address"] != "" {
		HttpListenAddress = c["http_listen_address"]
	}
	if c["key_path"] != "" {
		KeyPath = c["key_path"]
	}
	if c["repo_type"] != "" {
		RepoType = c["repo_type"]
	}
	if c["repo_location"] != "" {
		RepoLocation = c["repo_location"]
	}
	if c["key_auth_type"] != "" {
		KeyAuthType = c["key_auth_type"]
	}
	if c["key_auth_location"] != "" {
		KeyAuthLocation = c["key_auth_location"]
	}
	if c["repo_type"] != "" {
		RepoType = c["repo_type"]
	}
	if c["repo_location"] != "" {
		RepoLocation = c["repo_location"]
	}
	if c["repo_type"] != "" {
		RepoType = c["repo_type"]
	}
	if c["repo_location"] != "" {
		RepoLocation = c["repo_location"]
	}
	if c["repo_type"] != "" {
		RepoType = c["repo_type"]
	}
	if c["repo_location"] != "" {
		RepoLocation = c["repo_location"]
	}
	if c["token"] != "" {
		Token = c["token"]
	}
}