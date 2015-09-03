package templates

import "errors"

var (
	ApiUnavailableError = errors.New("nanobox api is unavailable")
	ApiUnavailable      = []byte(`:::::::::::::::: NANOBOX API UNAVAILABLE !!!

Weird... we are unable to connect to our own 
API. Give us a minute, and try again. 

:::::::::::::::::::::::::::::::::::::::::::

Disconnecting.`)
)
