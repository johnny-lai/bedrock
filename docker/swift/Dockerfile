FROM debian:jessie

RUN apt-get update && apt-get install -y libedit-dev wget

WORKDIR /root

RUN wget https://swift.org/builds/ubuntu1404/swift-2.2-SNAPSHOT-2015-12-01-b/swift-2.2-SNAPSHOT-2015-12-01-b-ubuntu14.04.tar.gz

RUN tar zxvf swift-2.2-SNAPSHOT-2015-12-01-b-ubuntu14.04.tar.gz && \
    cp -R swift-2.2-SNAPSHOT-2015-12-01-b-ubuntu14.04/* / && \
    rm -rf swift-2.2-SNAPSHOT-2015-12-01-b-ubuntu14.04 && \
    rm swift-2.2-SNAPSHOT-2015-12-01-b-ubuntu14.04.tar.gz

