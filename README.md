# YAAC
[![GitHub license](https://img.shields.io/github/license/DHBW-SE-2023/YAAC.svg)](https://github.com/DHBW-SE-2023/YAAC/blob/master/LICENSE) 
Full-Stack YAAC Prototype writen in Go.

# Installation
## Prerequisites
- Windows
- Docker
- WSL2

## Install the latest image
Pull the docker image:
```
docker pull ghcr.io/dhbw-se-2023/yaac:nightly
```

Run the docker image (this command must be executed in WSL2)
```
docker run -it --rm -e DISPLAY=$DISPLAY -v /tmp/.X11-unix:/tmp/.X11-unix ghcr.io/dhbw-se-2023/yaac:nightly
```

# Development Setup
## Prerequisites
Follow the guide on https://developer.fyne.io/started/ for your system. \
And follow the steps on https://gocv.io/getting-started/ for your system.
## Running
```shell
make run
```
## Packaging
```shell
make package
```
