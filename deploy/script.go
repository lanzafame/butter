package deploy

import (
)

type (
	script struct {}
)

func init() {
	PassRegister("script", script{})
}

func (s script) Initialize() error {
	return nil
}

func (s script) Run(stream io.Writer, id string) (error) {
	return nil
}
