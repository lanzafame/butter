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
	"fmt"
	"github.com/pagodabox/golang-mist"
	"github.com/pagodabox/na-ssh/templates"
	"github.com/pagodabox/nanobox-config"
	"io"
	"net/http"
	"time"
)

func Deploy(stream io.Writer, deploy string) error {
	done := make(chan bool, 0)
	disconnected := make(chan bool, 0)

	stream.Write(templates.Header)
	// connect up to mist
	client, err := mist.NewRemoteClient(config.Config["mistAddress"])
	if err != nil {
		stream.Write(templates.ApiUnavailable)
		return templates.ApiUnavailableError
	}
	defer client.Close()
	client.Subscribe([]string{"log", "deploy"})

	// TODO: connect up to noanobox mist (websocket)
	go func() {
		// subscribe to transaction updates? maybe show which transaction is running?
		transactions := []string{"deploy", "scale", "rename", "scale"}
		for _, name := range transactions {
			stream.Write([]byte(fmt.Sprintf("waiting on transaction:%-15v\r", name)))
			time.Sleep(time.Second * 2)
		}
		// clear the line
		stream.Write([]byte("                                     \r"))
		logs := []string{
			"Connecting to build server",
			"Starting build",
			"[√] submodules installed",
			"Build Sucessfull",
			"Uploading build",
			"Build uploaded",
			"Ordering servers",
			"Provisioning servers",
			"Adding servers to routing mesh",
			"Switching over routing mesh",
		}
		for _, log := range logs {
			client.Publish([]string{"log", "deploy"}, log)
			time.Sleep(time.Second * 1)
		}
		// signal that the current deploy is done
		close(done)
	}()

	// TODO: send deploy to nanobox
	http.Post("https://nanobox.io/apps/me/deploy/"+deploy, "application/json", nil)

	// listen for changes
outer:
	for {
		select {
		case msg, ok := <-client.Messages():
			if !ok {
				stream.Write(templates.ApiUnavailable)
				return templates.ApiUnavailableError
			}
			// write the log
			fmt.Println(msg)
			stream.Write([]byte(msg.Data.(string) + "\n"))
		case <-done:
			break outer
		case <-disconnected:
			stream.Write(templates.Disconnected)
			return templates.DisconnectedError
		}
	}

	stream.Write(templates.Footer)
	return nil
}
