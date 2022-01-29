FROM golang:1.17-alpine

WORKDIR /vestiarie-test-task

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -a -o bin/vestiarie-test-task main.go

EXPOSE 8080

CMD [ "/usr/bin/vestiarie-test-task" ]