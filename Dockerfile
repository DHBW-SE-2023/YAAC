FROM gocv/opencv:latest
# Debian based
MAINTAINER YAAC Team
# Update
RUN apt-get update && apt-get upgrade -y
# Get General Dependencies
RUN apt-get install -y git make
# Get Fyne Dependencies
RUN apt-get install -y gcc libgl1-mesa-dev xorg-dev mesa-utils
# Get gocv dependencies
#RUN apt-get
# Get Repo - WORKDIR is /go
RUN git clone https://github.com/DHBW-SE-2023/YAAC.git
#CMD ["go", "run", "src/gocv.io/x/gocv/cmd/version/main.go"]
WORKDIR "/go/YAAC"
RUN make build
CMD ["make", "run"]