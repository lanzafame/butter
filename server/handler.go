package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"github.com/nanopack/butter/repo"
)

type (
	Handler interface {
		Request(ssh.Channel, *ssh.Request) (bool, error)
	}

	handle struct {
		env map[string]string
	}
)

var (
	UnknownChannel    = errors.New("Unknown channel type")
	NotYetImplemented = errors.New("Not yet implemented")
	ShellDisabled     = errors.New("Shell access is not allowed")
	UnknownCommand    = errors.New("Unknown command")
	UnkownRequest     = errors.New("Unknown request")
	Done              = errors.New("Done")
)

func NewHandle(name string) (Handler, error) {
	switch name {
	case "session":
		handle := handle{
			env: make(map[string]string, 0),
		}
		return &handle, nil
	case "direct-tcpip":
		return nil, NotYetImplemented
	default:
		return nil, UnknownChannel
	}
}

func (handle *handle) Request(ch ssh.Channel, req *ssh.Request) (bool, error) {
	switch req.Type {
	case "pty-req":
		fallthrough
	case "shell":
		return true, ShellDisabled
	case "env":
		// we store these off??
		return false, nil
	case "exec":
		// it is prefixed with the length so we strip it off
		command := string(req.Payload[4:])

		// find the correct handler and run it
		for _, cmd := range repo.Commands() {
			if cmd.Match(command) {
				fmt.Println("found match", command)
				code, err := cmd.Run(command, ch)
				exitStatusBuffer := make([]byte, 4)
				binary.PutUvarint(exitStatusBuffer, uint64(code))
				fmt.Println("cmd finished", code, err)
				// purposefully ignoring the possible error
				ch.SendRequest("exit-status", false, exitStatusBuffer)
				return true, err
			}
		}

		return true, UnknownCommand
	}
	return true, UnkownRequest
}
