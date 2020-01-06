# go-libsrv

Library for all GO server project

## Requirement

need environment variable

SA : {google service account key}

## Git

clone source code to local.

```bash
git clone git@github.com:piyuo/go-lib-server.git
```

## Dev

write test file and using go extension

```bash
run test | debug test
```

or

using launch.json configuration

```json
{
 "name": "libsrv debug",
 "type": "go",
 "request": "launch",
 "mode": "auto",
 "program": "${workspaceFolder}/src/libsrv-debug"
}
```

set break point on main.go, trace into server.Start() to get libsrv server.go

use curl to debug

```bash
curl http://127.0.0.1:2999
```

## Test

unit test using go test

```bash
go test ./...
```
