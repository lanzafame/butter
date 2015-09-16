## Nanobox SSH

SSH server that handles authenticating users, forwards commands to git, and can be used to pull specific files from the git repo.

## Routes

| Route | Description | Payload | Response |
| --- | --- | --- | --- |
| `/files?ref={ref}` | Show the names of all the files at the specific ref, or MASTER | nil | `{file contents}` |
| `/files/{file}?ref={ref}` | Get the content of the file at the specific ref, or MASTER | nil | `{file contents}` |
| `/branches` | Get the names of all branches pushed | nil | `["master"]` |
| `/commits` | Get a list of all the commits | nil | `[{"id":"sha","message":"this is a message","author":"me"}]` |
| `/commits/{commit}` | Get details about a specific commit | nil | `[{"id":"sha","author":"me","message":"this is a message","author_date":"jan","author_email":"me@me.com"}]` |

## TODO

- add a handler for tunnels?