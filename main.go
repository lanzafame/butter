// -*- mode: go; tab-width: 2; indent-tabs-mode: 1; st-rulers: [70] -*-
// vim: ts=4 sw=4 ft=lua noet
//--------------------------------------------------------------------
// @author Daniel Barney <daniel@nanobox.io>
// Copyright (C) Pagoda Box, Inc - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly
// prohibited. Proprietary and confidential
//
// @doc
//
// @end
// Created :   1 September 2015 by Daniel Barney <daniel@nanobox.io>
//--------------------------------------------------------------------
package main

import (
	"bitbucket.org/nanobox/na-api"
	"bitbucket.org/nanobox/na-ssh/commands"
	"bitbucket.org/nanobox/na-ssh/git"
	"bitbucket.org/nanobox/na-ssh/handler"
	"bitbucket.org/nanobox/na-ssh/nanobox"
	"bitbucket.org/nanobox/na-ssh/routes"
	"bitbucket.org/nanobox/na-ssh/server"
	nanoboxConfig "github.com/pagodabox/nanobox-config"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
)

var (
	config    map[string]string
	sshServer *ssh.ServerConfig
)

func init() {
	defaults := map[string]string{
		"listenAddress":     ":2222",
		"httpListenAddress": ":2222",
		"keyPath":           "./host_key",
		"gitRepo":           "./testing",
		"mistAddress":       "127.0.0.1:1445",
	}

	nanoboxConfig.Load(defaults, "")
	config = nanoboxConfig.Config

	hostPrivateKey, err := ioutil.ReadFile(config["keyPath"])
	if err != nil {
		panic(err)
	}

	hostPrivateKeySigner, err := ssh.ParsePrivateKey(hostPrivateKey)
	if err != nil {
		panic(err)
	}

	sshServer = &ssh.ServerConfig{
		PublicKeyCallback: nanobox.Authenticate,
	}

	sshServer.AddHostKey(hostPrivateKeySigner)

	// check if the git repo is already set up
	os.Mkdir(config["gitRepo"], 0700)
	for _, name := range []string{"live.git", "staging.git"} {
		err = git.Init(config["gitRepo"] + "/" + name)
		if err != nil {
			panic(err)
		}
	}

	// add in our custom commands that the ssh server can respond to
	handler.Commands = append(handler.Commands, commands.Push{})
	handler.Commands = append(handler.Commands, commands.Pull{})
}

func main() {
	server, err := server.StartServer(config["listenAddress"], sshServer)
	if err != nil {
		panic(err)
	}
	defer server.Close()

	routes.Init()

	api.Start(config["httpListenAddress"])
}
