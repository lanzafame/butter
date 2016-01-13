[![butter logo](http://nano-assets.gopagoda.io/readme-headers/butter.png)](http://nanobox.io/open-source#butter)
[![Build Status](https://travis-ci.org/nanopack/butter.svg)](https://travis-ci.org/nanopack/butter)

## Butter

A small, version controll based deployment service with pluggable authentication and deployment strategies.

### Status
Experimental/Unstable/Incomplete

## Routes

| Route | Description | Payload | Response |
| --- | --- | --- | --- |
| `/files?ref={ref}` | Show the names of all the files at the specific ref, or MASTER | nil | `{file contents}` |
| `/files/{file}?ref={ref}` | Get the content of the file at the specific ref, or MASTER | nil | `{file contents}` |
| `/branches` | Get the names of all branches pushed | nil | `["master"]` |
| `/commits` | Get a list of all the commits | nil | `[{"id":"sha","message":"this is a message","author":"me"}]` |
| `/commits/{commit}` | Get details about a specific commit | nil | `[{"id":"sha","author":"me","message":"this is a message","author_date":"jan","author_email":"me@me.com"}]` |

[![butter logo](http://nano-assets.gopagoda.io/open-src/nanobox-open-src.png)](http://nanobox.io/open-source)


## TODO
build a cli
Write tests