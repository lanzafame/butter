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
	"github.com/pagodabox/nanobox-config"
	"golang.org/x/crypto/ssh"
	"os"
	"strings"
)

var config map[string]string

func init() {
	defaults := map[string]string{
		"listenAddress": ":2222",
	}

	config.Load(defaults, "")
	config = config.Config
}

func main() {

}
