FROM golang:1.3

ADD . /go/src/github.com/flowcker/flowcker-stdlib

RUN go get github.com/flowcker/flowcker-stdlib/cmd

RUN go install github.com/flowcker/flowcker-stdlib/cmd

ENV PORT 3000

ENTRYPOINT ["/go/bin/stdlib"]

EXPOSE 3000
