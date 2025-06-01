FROM ubuntu:20.04

RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y \
    build-essential \
    libglib2.0-dev \
    flex \
    libgtk2.0-dev \
    libglade2-dev \
    libreadline-dev \
    guile-2.2-dev \
    intltool \
    pkg-config \
    gettext \
    && rm -rf /var/lib/apt/lists/*

COPY mdk-1.3.0.tar.gz /tmp/
WORKDIR /tmp
RUN tar xfvz mdk-1.3.0.tar.gz && \
    cd mdk-1.3.0 && \
    ./configure && \
    make && \
    make install

WORKDIR /workspace
VOLUME ["/workspace"]
CMD ["/bin/bash"]