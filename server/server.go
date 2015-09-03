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
package server

import (
	"fmt"
	"github.com/pagodabox/na-ssh/handler"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
)

func StartServer(address string, sshServer *ssh.ServerConfig) (io.Closer, error) {
	serverSocket, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := serverSocket.Accept()
			if err != nil {
				return
			}
			go handleConnection(conn, sshServer)
		}
	}()
	return serverSocket, nil
}

func handleConnection(conn net.Conn, sshServer *ssh.ServerConfig) {

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, sshServer)
	if err != nil {
		fmt.Println("got an error", err)
		return
	}
	fmt.Println("connection was established")

	defer sshConn.Close()

	go func() {
		for req := range reqs {
			fmt.Printf("got an out-of-band request '%v'\n", req.Type)
			// got an out-of-band request
			// we ignore these? or what?
		}
	}()
	for chanRequest := range chans {
		handleChannel(chanRequest)
	}
}

func handleChannel(chanRequest ssh.NewChannel) {
	handle, err := handler.NewHandle(chanRequest.ChannelType())
	if err != nil {
		fmt.Println("wrong handler", err)
		chanRequest.Reject(ssh.UnknownChannelType, err.Error())
		return
	}

	ch, reqs, err := chanRequest.Accept()
	if err != nil {
		fmt.Println("fail to accept channel request", err)
		return
	}

	defer ch.Close()

	for req := range reqs {
		done, err := handle.Request(ch, req)
		if err != nil {
			fmt.Println("request errored out", err)
			_, err := ch.Write([]byte(fmt.Sprintf("%v\r\n", err)))
			fmt.Println(err)
		}
		if done {
			return
		}
	}
}
