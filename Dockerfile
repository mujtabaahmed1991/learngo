FROM golang:1.9.3
ADD . /go/src/github.com/golearn
WORKDIR /go/src/github.com/golearn

RUN go get -u github.com/unrolled/render

RUN go install github.com/golearn

ENTRYPOINT /go/bin/golearn

EXPOSE 8080
