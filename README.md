# LaraGo Installer Tool

It was made to help you start new LaraGo project as easy as possible.

1. Download [LaraGo Boilerplate](https://github.com/lara-go/boilerplate) in zip.
1. Unpack it to the directory you provided (it should not exist).
1. Update import paths.

## Installation

```bash
go get -u github.com/lara-go/installer
```

## Running

```bash
$GOPATH/bin/installer install my/application/path
```

Result will be:
```
New project was installed in /Users/user/go/src/my/application/path

Next steps:

  1. Install Glide tool (https://glide.sh)
  2. Install dependencies using Glide

Check if everyting ok:

  $ go run cmd/app/main.go -r ./ env
```
