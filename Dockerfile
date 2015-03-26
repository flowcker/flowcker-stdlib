FROM golang:1.3

ADD . /go/src/github.com/flowcker/flowcker-stdlib

RUN go get -v github.com/flowcker/flowcker-stdlib/cmd/stdlib

RUN go install -v github.com/flowcker/flowcker-stdlib/cmd/stdlib

ENV PORT 3000

ENTRYPOINT ["/go/bin/stdlib"]

EXPOSE 3000
