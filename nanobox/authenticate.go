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
package nanobox

import (
	"github.com/pagodabox/na-ssh/templates"
	"golang.org/x/crypto/ssh"
	"io"
)

func Authenticate(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {

	// TODO: We connect to nanobox and authenticate the user for the app

	// store off what we need here
	// return &ssh.Permissions{Extensions: map[string]string{"user_id": user.Id}}, nil
	return nil, nil
}

func Deploy(stream io.Writer, deploy string) error {
	stream.Write(templates.Header)
	stream.Write(templates.Footer)
	return nil
}
