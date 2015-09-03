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
package git

import (
	"fmt"
	"github.com/pagodabox/nanobox-config"
	"io"
	"os"
	"os/exec"
)

func Init(path string) error {
	os.Mkdir(path, 0777)
	if _, err := os.Stat(path + "/info"); os.IsNotExist(err) {
		cmd := exec.Command("git", "init", "--bare")
		cmd.Dir = path
		return cmd.Run()
	}
	return nil
}

func Shell(duplex io.ReadWriter, errStream io.Writer, command string) (uint64, error) {
	cmd := exec.Command("git", "shell", "-c", command)
	cmd.Dir = config.Config["gitRepo"]

	cmd.Stdout = duplex
	cmd.Stderr = errStream
	cmd.Stdin = duplex

	err := cmd.Run()
	fmt.Println(err)

	if err != nil {
		// should return the actual exit code
		return 1, err
	}

	return 0, nil
}
