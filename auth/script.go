package auth

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"os/exec"

	"github.com/nanopack/butter/config"
)

type (
	ScriptPassAuth struct{}
	ScriptKeyAuth  struct{}
)

func init() {
	PassRegister("script", ScriptPassAuth{})
	KeyRegister("script", ScriptKeyAuth{})
}

func (s ScriptPassAuth) Initialize() error {
	file, err := os.Open(config.PassAuthLocation)
	if err != nil {
		return fmt.Errorf("ScriptPassAuth: %+v", err)
	}
	fi, err := file.Stat()
	if err != nil {
		return fmt.Errorf("ScriptPassAuth: %+v", err)
	}
	if fi.IsDir() {
		return fmt.Errorf("file given is a directory")
	}
	return nil
}

func (s ScriptKeyAuth) Initialize() error {
	file, err := os.Open(config.KeyAuthLocation)
	if err != nil {
		return fmt.Errorf("ScriptKeyAuth: %+v", err)
	}
	fi, err := file.Stat()
	if err != nil {
		return fmt.Errorf("ScriptKeyAuth: %+v", err)
	}
	if fi.IsDir() {
		return fmt.Errorf("file given is a directory")
	}
	return nil
}

func (s ScriptPassAuth) Auth(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	cmd := exec.Command(config.PassAuthLocation, conn.User(), conn.RemoteAddr().String())
	passReader := bytes.NewReader(password)
	cmd.Stdin = passReader
	output, err := cmd.CombinedOutput()
	if err != nil {
		config.Log.Error("password authentication: %s\n%v", output, err)
		return nil, err
	}

	// nil permissions is success?
	return nil, nil
}

func (s ScriptKeyAuth) Auth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	cmd := exec.Command(config.KeyAuthLocation, conn.User(), conn.RemoteAddr().String())
	k := fmt.Sprintf("%s\n",base64.StdEncoding.EncodeToString(key.Marshal()))
	keyReader := bytes.NewReader([]byte(k))
	cmd.Stdin = keyReader
	output, err := cmd.CombinedOutput()
	if err != nil {
		config.Log.Error("key authentication: %s\n%v", output, err)
		return nil, err
	}

	// nil permissions is success?
	return nil, nil
}
