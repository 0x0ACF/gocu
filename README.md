
# Gocu

## Overview
Gocu is [curl](https://curl.se/) reimagined, a CLI HTTP client written in Go focused on simplicity and ease of use. It only implements a small subset of curl options, adding some features such as **placeholders** for variables.

## Installation
Gocu can be installed in your GOPATH using the following command (Go version 1.24.2 or higher is required)
```
go install github.com/0x0ACF/gocu@latest
```

## Usage
```
Usage:
  gocu [flags]
  gocu [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  vars        Manage the variables used as placeholders

Flags:
  -d, --data string          Data to send in the body of the request
  -H, --header stringArray   Header to add to the request
  -h, --help                 Help for gocu
  -X, --request string       HTTP method to use (default "GET")
```

Variables are one of the main features of gocu. You can manage variables with the **vars** subcommand, and use them when doing requests wrapping their names between brackets ({{baseAddress}}, for example).
All variables are saved in a file called **gocu.cache** in `os.UserCacheDir()`.
```
Usage:
  gocu vars [flags]
  gocu vars [command]

Available Commands:
  add         Adds a new variable
  clear       Removes all saved variables
  get         Gets a variable value
  ls          Lists all saved variables
  mod         Modifies a variable value
  rm          Removes a variable

Flags:
  -h, --help   help for vars
```

## Examples
- Perform a GET request
```
gocu http://example.org
```

- Perform a POST request with a JSON payload
```
gocu http://example.org -X POST -d '{"key": "value"}'
```

- Perform a PUT request with a JSON payload and two headers
```
gocu http://example.org -X PUT -d '{"key": "value"}' -H "Authorization: Bearer base64token" -H "Accept: application/json"
```

- Add two variables and use them
```
gocu vars add baseAddress http://example.org
gocu vars add token mytoken
gocu {{baseAddress}} -H "Authorization: {{token}}"
```
