package auth

import (
	"github.com/nanopack/butter/config"
	"golang.org/x/crypto/ssh"
)

type (
	KeyAuther interface {
		Initialize() error
		Auth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error)
	}

	PassAuther interface {
		Initialize() error
		Auth(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error)
	}
)

var (
	availableKeyAuthers = map[string]KeyAuther{}
	defaultKeyAuther    KeyAuther

	availablePassAuthers = map[string]PassAuther{}
	defaultPassAuther    PassAuther
)

func KeyRegister(name string, k KeyAuther) {
	availableKeyAuthers[name] = k
}

func PassRegister(name string, p PassAuther) {
	availablePassAuthers[name] = p
}

func Setup() error {
	keyauth, ok := availableKeyAuthers[config.KeyAuthType]
	if ok {
		config.Log.Info("setting up key auth(%s) location: %s", config.KeyAuthType, config.KeyAuthLocation)
		defaultKeyAuther = keyauth
		if err := keyauth.Initialize(); err != nil {
			return err
		}
	}
	passauth, ok := availablePassAuthers[config.PassAuthType]
	if ok {
		config.Log.Info("setting up pass auth(%s) location: %s", config.PassAuthType, config.PassAuthLocation)
		defaultPassAuther = passauth
		if err := passauth.Initialize(); err != nil {
			return err
		}
	}
	return nil
}

func KeyAuth() func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	if defaultKeyAuther == nil {
		return nil
	}
	return defaultKeyAuther.Auth
}

func PassAuth() func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	if defaultPassAuther == nil {
		return nil
	}
	return defaultPassAuther.Auth
}
