# libsrv

Library for all server project

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
  "name": "libsrv debug test",
  "type": "go",
  "request": "launch",
  "mode": "test",
  "program": "${workspaceFolder}/src/libsrv"
}
```

## Test

unit test using go test

```bash
go test ./...
```
