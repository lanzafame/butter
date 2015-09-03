package templates

var Unauthorized = []byte(`::::::::::::::::: AUTHORIZATION FAILURE !!!

We're sorry, the SSH key provided is not 
authorized for this repo.

Please verify that:

1- Your account is verified
2- Your SSH key is provided and correct
3- You own, or have been given collaboration
   rights to this repository.

:::::::::::::::::::::::::::::::::::::::::::

Disconnecting.`)
