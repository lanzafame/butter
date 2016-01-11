package server

import (
	"fmt"
	"github.com/nanopack/butter/config"
	"github.com/nanopack/butter/auth"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
)

var sshServer *ssh.ServerConfig

func StartServer() (io.Closer, error) {
	hostPrivateKey, err := ioutil.ReadFile(config.KeyPath)
	if err != nil {
		panic(err)
	}

	hostPrivateKeySigner, err := ssh.ParsePrivateKey(hostPrivateKey)
	if err != nil {
		panic(err)
	}

	sshServer = &ssh.ServerConfig{
		PasswordCallback: auth.PassAuth(),
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
	return serverSocket, nil
}

func handleConnection(conn net.Conn) {

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
	handle, err := NewHandle(chanRequest.ChannelType())
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
