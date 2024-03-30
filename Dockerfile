FROM docker.io/gocv/opencv:4.8.1
# Debian based
LABEL org.opencontainers.image.authors="YAAC Team"
# Update
RUN apt-get update \
    && apt-get upgrade -y
RUN apt-get update \
    && apt-get install --no-install-recommends -y software-properties-common \
    && apt-get clean
# Get General Dependencies
RUN apt-get update \
    && apt-get install --no-install-recommends -y git make mesa-utils libglfw3 \
    && apt-get clean
# Get Fyne Dependencies
RUN apt-get update \
    && apt-get install --no-install-recommends -y gcc libgl1-mesa-dev xorg-dev \
    && apt-get clean
# Install tesseract deps
RUN apt-get update \
    && apt-get install --no-install-recommends -y tesseract-ocr tesseract-ocr-osd tesseract-ocr-eng tesseract-ocr-deu libtesseract4 libtesseract-dev \
    && apt-get clean
# Get OpenGL Deps (LLVMPipe)
RUN apt-get update \
    && apt-get install --no-install-recommends -y meson llvm libx11-dev libxrandr-dev libxi-dev \
    && apt-get clean
# Try to use software rendereing
#ENV LIBGL_ALWAYS_SOFTWARE=1
# Get Repo - WORKDIR is /go
WORKDIR "/go"
RUN git clone https://github.com/DHBW-SE-2023/YAAC.git
WORKDIR "/go/YAAC"
RUN make build
CMD ["make", "run"]