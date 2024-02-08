FROM ubuntu:latest
MAINTAINER YAAC Team
RUN apt-get update && apt-get upgrade -y
CMD ["ls"]