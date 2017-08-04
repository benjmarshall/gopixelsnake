# Go Pixel Snake

Go Pixel Snake is a weekend project to create a Snake game in Go, using the [Pixel game library](https://github.com/faiface/pixel).

![Go Pixel Snake](/resources/screenshot.png)

## Install
### Build from source
The Pixel game engine uses an Opengl back-end which requires [platform specific dependencies](https://github.com/faiface/pixel#requirements) for compilation. You will have to install these dependencies on your system first.

Go Pixel Snake uses the [go dep](https://github.com/golang/dep) tool to manage its own dependencies.
```
go get github.com/golang/dep
go get github.com/benjmarshall/gopixelsnake
cd $GOPATH/src/github.com/benjmarshall/gopixelsnake
dep ensure
go install
```
### Pre-built
Alternatively download one of the pre-built binaries from the releases page.

### Bugs
There are probably many bugs in here. If you spot something major please submit an issue.
