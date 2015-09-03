package templates

import "errors"

var (
	DisconnectedError = errors.New("disconnected from the api")
	Disconnected      = []byte(`:::::::::::: DEPLOY STREAM DISCONNECTED !!!

Oh snap the deploy stream just disconnected. 
No worries, you can visit the dashboard to
view the complete output stream.

:::::::::::::::::::::::::::::::::::::::::::

Disconnecting.`)
)
