FROM golang:1.14.1

WORKDIR /http_server

COPY main.go .

RUN go build

RUN chmod +x /http_server/http_server

ENTRYPOINT [ "/http_server/http_server" ]

