FROM docker.io/gocv/opencv:latest
# Debian based
MAINTAINER YAAC Team
# Update
RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y software-properties-common
RUN apt-get update
# Get General Dependencies
RUN apt-get install -y git make mesa-utils libglfw3
# Get Fyne Dependencies
RUN apt-get install -y gcc libgl1-mesa-dev xorg-dev
# Get OpenGL Deps (LLVMPipe)
RUN apt-get -y install meson llvm libx11-dev libxrandr-dev libxi-dev
RUN apt-get clean
ENV LIBGL_ALWAYS_SOFTWARE=1
#WORKDIR "/"
#RUN mkdir build
#WORKDIR "/build"
#RUN meson -D glx=xlib -D gallium-drivers=swrast
#RUN ninja
#RUN export LD_LIBRARY_PATH="lib/gallium/libGL.so"
# Get Repo - WORKDIR is /go
WORKDIR "/go"
RUN git clone https://github.com/DHBW-SE-2023/YAAC.git
#CMD ["go", "run", "src/gocv.io/x/gocv/cmd/version/main.go"]
WORKDIR "/go/YAAC"
RUN make build
CMD ["make", "run"]