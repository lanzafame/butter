package server

import (
	"fmt"
	"github.com/nanopack/butter/auth"
	"github.com/nanopack/butter/config"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
)

var sshServer *ssh.ServerConfig

func StartServer() (io.Closer, error) {
	hostPrivateKey, err := ioutil.ReadFile(config.KeyPath)
	if err != nil {
		return nil, err
	}

	hostPrivateKeySigner, err := ssh.ParsePrivateKey(hostPrivateKey)
	if err != nil {
		return nil, err
	}

	sshServer = &ssh.ServerConfig{
		PasswordCallback:  auth.PassAuth(),
		PublicKeyCallback: auth.KeyAuth(),
	}

	sshServer.AddHostKey(hostPrivateKeySigner)

	serverSocket, err := net.Listen("tcp", config.SshListenAddress)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := serverSocket.Accept()
			if err != nil {
				return
			}
			go handleConnection(conn)
		}
	}()
	config.Log.Info("SSH listening on %s", config.SshListenAddress)
	return serverSocket, nil
}

func handleConnection(conn net.Conn) {

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, sshServer)
	if err != nil {
		config.Log.Debug("got an error", err)
		return
	}
	config.Log.Debug("connection was established")

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
	handle, err := NewHandle(chanRequest.ChannelType())
	if err != nil {
		config.Log.Debug("wrong handler %v", err)
		chanRequest.Reject(ssh.UnknownChannelType, err.Error())
		return
	}

	ch, reqs, err := chanRequest.Accept()
	if err != nil {
		config.Log.Debug("fail to accept channel request %v", err)
		return
	}

	defer ch.Close()

	for req := range reqs {
		done, err := handle.Request(ch, req)
		if err != nil {
			config.Log.Debug("request errored out %v", err)
			_, err := ch.Write([]byte(fmt.Sprintf("%v\r\n", err)))
			if err != nil {
				config.Log.Debug(err.Error())
			}
		}
		if done {
			return
		}
	}
}
