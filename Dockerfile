FROM golang:1.17-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o vestiarie-test-task .

WORKDIR /dist
RUN cp /build/vestiarie-test-task .

EXPOSE 8080
