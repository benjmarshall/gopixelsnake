# Go Pixel Snake

Go Pixel Snake is a weekend project to create a Snake game in Go, using the [Pixel game engine](https://github.com/faiface/pixel).

![Go Pixel Snake](/resources/screenshot.png)

## Install
### Build from source
The Pixel game engine uses an opengl backend which requires [platform specific dependencies](https://github.com/faiface/pixel#requirements) for compilation. You will have to install these dependencies on your system first.

Go Pixel Snake uses the [go dep](https://github.com/golang/dep) tool to manage it's own dependencies.
```
go get github.com/golang/dep
go get github.com/benjmarshall/gopixelsnake
cd $GOPATH/src/github.con/benjmarshall/gopixelsnake
dep ensure
go install
```
### Pre-built
Alternatively download one of the pre-built binaries from the releases page.
