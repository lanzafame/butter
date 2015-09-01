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
	"fmt"
	nanoboxConfig "github.com/pagodabox/nanobox-config"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
)

var (
	config    map[string]string
	sshServer *ssh.ServerConfig
)

func init() {
	defaults := map[string]string{
		"listenAddress": ":2222",
		"keyPath":       "./host_key",
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
		PublicKeyCallback: keyAuth,
	}

	sshServer.AddHostKey(hostPrivateKeySigner)
}

func main() {
	socket, err := net.Listen("tcp", config["listenAddress"])
	if err != nil {
		panic(err)
	}

	for {
		conn, err := socket.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	// upgrade to ssh connection
	fmt.Println("new connection")

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
	switch chanRequest.ChannelType() {
	case "direct-tcpip":
		// this is for port forwarding
		chanRequest.Reject(ssh.UnknownChannelType, "Not Yet Implemented")
	case "session":
		ch, reqs, err := chanRequest.Accept()
		if err != nil {
			fmt.Println("fail to accept channel request", err)
			return
		}
		for req := range reqs {
			handleChannelRequest(ch, req)
		}
	default:
		chanRequest.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", chanRequest.ChannelType()))
	}
}

func handleChannelRequest(ch ssh.Channel, req *ssh.Request) {
	switch req.Type {
	case "pty-req":
		fallthrough
	case "shell":
		ch.Write([]byte("shell access is not allowed\r\n"))
		ch.Close()
	case "env":
		// do we store these off?
	case "exec":

		// it is prefixed with the length so we strip it off
		command := string(req.Payload[4:])

		switch {
		case strings.HasPrefix(command, "git-receive-pack "):
			fallthrough
		case strings.HasPrefix(command, "git-upload-pack "):
			fmt.Printf("got command '%v'\n", command)
			cmd := exec.Command("git", "shell", "-c", command)

			outPipe, err := cmd.StdoutPipe()
			if err != nil {
				fmt.Println("got an error", err)
				return
			}
			errPipe, err := cmd.StderrPipe()
			if err != nil {
				fmt.Println("got an error", err)
				return
			}
			inPipe, err := cmd.StdinPipe()
			if err != nil {
				fmt.Println("got an error", err)
				return
			}
			if err := cmd.Start(); err != nil {
				fmt.Println("got an error", err)
			}

			go io.Copy(ch, outPipe)
			go io.Copy(ch, errPipe)
			go io.Copy(inPipe, ch)

			err = cmd.Wait()
			if err != nil {
				fmt.Println("got an error", err)
			}
			fmt.Println("closing connection")
			ch.Close()
		case strings.HasPrefix(command, "tunnel"):
			ch.Write([]byte("establishing tunnel! (NOT YET IMPLEMENTED)\r\n"))
			ch.Close()
		default:
			ch.Write([]byte(fmt.Sprintf("unsupported command: '%v'\r\n", []byte(command))))
			ch.Close()
		}

	default:
		ch.Write([]byte(fmt.Sprintf("request type '%v' is not implemented\r\npayload: %v\r\n", req.Type, string(req.Payload))))
		// ch.Close()
	}
}

func keyAuth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {

	// TODO: We connect to nanobox and authenticate the user for the app

	// store off what we need here
	// return &ssh.Permissions{Extensions: map[string]string{"user_id": user.Id}}, nil
	return nil, nil
}
