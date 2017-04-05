# LaraGo Installer Tool

It was made to help you start new LaraGo project as easy as possible.

1. Download LaraGo Boilerplate in zip.
1. Unpack it to the directory you need.
1. Update import paths.

## Installation

```bash
$ go install https://github.com/lara-go/installer
```

## Running

```bash
$ larago install my/application/path
```

Result will be:
```
New project was installed in /Users/user/go/src/my/application/path

Next steps:

  1. Install Glide tool (https://glide.sh)
  2. Install dependencies using Glide

To check run:

  $ go run cmd/app/main.go -r ./ env
```
