FROM golang:1.9.3
ADD ./code /go/src/github.com/exercise_gitlab
WORKDIR /go/src/github.com/exercise_gitlab

RUN go get -u github.com/unrolled/render

RUN go install github.com/exercise_gitlab/code

ENTRYPOINT /go/bin/exercise_gitlab

EXPOSE 8080
