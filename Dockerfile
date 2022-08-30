FROM golang:1.18

WORKDIR /go/src/terraform-provider-dmsnitch
COPY . .

ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go install -v ./...

VOLUME /go/src/terraform-provider-dmsnitch
