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
// Created :   2 September 2015 by Daniel Barney <daniel@nanobox.io>
//--------------------------------------------------------------------
package commands

import (
	"github.com/pagodabox/na-ssh/git"
	"github.com/pagodabox/na-ssh/nanobox"
	"github.com/pagodabox/na-ssh/templates"
	"github.com/pagodabox/nanobox-config"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"strings"
)

type (
	Push struct{}
	Pull struct{}
)

func (push Push) Match(command string) bool {
	return strings.HasPrefix(command, "git-receive-pack ")
}

func (push Push) Run(command string, ch ssh.Channel) (uint64, error) {
	//TODO make "master" be dynamic
	originalCommit := getCommit("live.git", "master")
	code, err := git.Shell(ch, ch.Stderr(), command)
	if err == nil {
		stream := ch.Stderr()

		newCommit := getCommit("live.git", "master")
		if newCommit == originalCommit {
			stream.Write(templates.NoChanges)
			return code, nil
		}

		err := nanobox.Deploy(stream, newCommit)

		switch {
		case err == templates.ApiUnavailableError:

			// write the original commit back to the master file, so that we
			// can trigger a deploy again without needing new code to be
			// pushed
			name := headName("live.git", "master")
			ioutil.WriteFile(name, []byte(originalCommit+"\n"), 0600)
			fallthrough

		case err != nil:

			// we return nil because we have already sent the error across
			// in the nanobox.Deploy function
			return 1, nil
		}
	}
	return code, err
}

func (pull Pull) Match(command string) bool {
	return strings.HasPrefix(command, "git-send-pack ")
}

func (pull Pull) Run(command string, ch ssh.Channel) (uint64, error) {
	return git.Shell(ch, ch.Stderr(), command)
}
func headName(repo, name string) string {
	return config.Config["gitRepo"] + "/" + repo + "/refs/heads/" + name
}
func getCommit(repo, name string) string {
	file := headName(repo, name)
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}
	return strings.TrimRight(string(bytes), "\n\r")
}
