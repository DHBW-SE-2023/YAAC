start /Min cmd /C wsl -- docker run -it --rm -e DISPLAY=$DISPLAY -v /tmp/.X11-unix:/tmp/.X11-unix -v $HOME:/host -v /mnt/c:/winc ghcr.io/dhbw-se-2023/yaac:master
