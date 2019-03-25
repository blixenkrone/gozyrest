FROM golang:alpine3.9 AS build
RUN apk update && apk upgrade && apk add --no-cache bash git openssh
LABEL version 1.0

ENV SRC_DIR=/go/blix/mongo
COPY . ${SRC_DIR}
WORKDIR ${SRC_DIR}
RUN go build -o mongo
ENTRYPOINT [ "./mongo" ]

EXPOSE 80
EXPOSE 8085