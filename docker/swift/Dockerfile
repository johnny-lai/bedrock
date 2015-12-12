FROM ubuntu:14.04

RUN apt-get update && \
  apt-get install -y wget \
    clang-3.6 git cmake ninja-build \
    uuid-dev libicu-dev icu-devtools libbsd-dev libedit-dev libxml2-dev \
    libsqlite3-dev swig libpython-dev libncurses5-dev pkg-config && \
  update-alternatives --install /usr/bin/clang clang /usr/bin/clang-3.6 100 && \
  update-alternatives --install /usr/bin/clang++ clang++ /usr/bin/clang++-3.6 100

WORKDIR /root

RUN wget https://swift.org/builds/ubuntu1404/swift-2.2-SNAPSHOT-2015-12-01-b/swift-2.2-SNAPSHOT-2015-12-01-b-ubuntu14.04.tar.gz

RUN tar zxvf swift-2.2-SNAPSHOT-2015-12-01-b-ubuntu14.04.tar.gz && \
    cp -R swift-2.2-SNAPSHOT-2015-12-01-b-ubuntu14.04/* / && \
    rm -rf swift-2.2-SNAPSHOT-2015-12-01-b-ubuntu14.04 && \
    rm swift-2.2-SNAPSHOT-2015-12-01-b-ubuntu14.04.tar.gz

