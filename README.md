# libsrv

Library for all GO server project

## Requirement

we separate source code from gopath so we don't need worried lot of pkg file sync over dropbox

Environment Variables

PATH="$HOME/go/bin":"$GOPATH/bin"
GO111MODULE=on
GOPATH="/Users/cc/gopath"
NAME="dev"
REGION="US"
BRANCH="master"

## Git

clone source code to local.

```bash
git clone git@github.com:piyuo/libsrv.git
```

## Test

unit test using go test

```bash
go test ./... -parallel 16
```

## Update go.mod

To upgrade all dependencies at once for a given module, just run the following from the root directory of your module

This upgrades to the latest or minor patch release

```bash
go get -u ./...
```

## use gopls

```bash
go get golang.org/x/tools/gopls@latest
```

### Dev

go clean cache

```bash
go clean -cache -modcache -i -r
```
