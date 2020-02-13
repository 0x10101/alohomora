FROM golang:alpine AS build-env

ENV SRC_DIR $GOPATH/src/github.com/steps0x29a/alohomora
ENV VERSION 0.5
RUN apk add --no-cache --update bash libpcap-dev git build-base

WORKDIR $SRC_DIR
ADD . $SRC_DIR
RUN go get -u github.com/google/gopacket
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/gosuri/uitable
RUN go get -u github.com/satori/go.uuid
RUN GOOS="linux" GOARCH="amd64" go build -v -ldflags "-w -s" -o alohomora main.go

FROM alpine
RUN apk add --no-cache --update bash libpcap
RUN mkdir /alohomora
RUN mkdir /data
VOLUME ["/data"]
COPY --from=build-env /go/src/github.com/steps0x29a/alohomora/alohomora /alohomora
WORKDIR /alohomora

EXPOSE 29100 29101
ENTRYPOINT ["/alohomora/alohomora", "--server", "-v"]
CMD ["--ip", "127.0.0.1"]