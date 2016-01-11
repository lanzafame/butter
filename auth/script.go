package auth

import (
	"golang.org/x/crypto/ssh"
)

type (
	ScriptPassAuth struct {}
	ScriptKeyAuth struct {}
)

func init() {
	PassRegister("script", ScriptPassAuth{})
	KeyRegister("script", ScriptKeyAuth{})
}

func (s ScriptPassAuth) Initialize() error {
	return nil
}

func (s ScriptKeyAuth) Initialize() error {
	return nil
}

func (s ScriptPassAuth) Auth(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	return nil, nil
}
func (s ScriptKeyAuth) Auth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	return nil, nil
}