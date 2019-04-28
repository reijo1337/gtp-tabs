FROM golang:latest

COPY . /go/src/gtp-tabs/GtpStorage/

RUN apt update && \
apt install unzip && \
curl -OL https://github.com/google/protobuf/releases/download/v3.2.0/protoc-3.2.0-linux-x86_64.zip && \
unzip protoc-3.2.0-linux-x86_64.zip -d protoc3 && \
mv protoc3/bin/* /usr/local/bin/ && \
mv protoc3/include/* /usr/local/include/

SHELL ["/bin/bash", "-c"]

RUN go get -u github.com/golang/protobuf/{proto,protoc-gen-go} && \
go get -u google.golang.org/grpc && \
# go get github.com/DATA-DOG/go-sqlmock \
# go get -u golang.org/x/net/context && \
# go get github.com/dgrijalva/jwt-go && \
go get github.com/lib/pq

WORKDIR /go/src/gtp-tabs/GtpStorage/protocol

RUN protoc --go_out=plugins=grpc:. *.proto

WORKDIR /go/src/gtp-tabs/GtpStorage

RUN go install gtp-tabs/GtpStorage

CMD [ "GtpStorage" ]