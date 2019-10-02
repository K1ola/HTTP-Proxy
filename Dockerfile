FROM ubuntu:18.04

#undone

RUN apt-get update && apt-get install -y wget

ENV GOVERSION 1.13
USER root
RUN wget https://storage.googleapis.com/golang/go$GOVERSION.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go$GOVERSION.linux-amd64.tar.gz
ENV GOROOT /usr/local/go
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH

COPY . .

CMD ./bin/proxy & ./bin/repeater