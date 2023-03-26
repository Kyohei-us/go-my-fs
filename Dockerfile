# syntax=docker/dockerfile:1

FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
RUN go mod download

COPY *.go ./
# COPY .env file if exists
# COPY .env ./
COPY *.html ./
# instead of VOLUME, create volume using command
# VOLUME /app/myfsfolder

RUN go build -o /go-my-fs

EXPOSE 8080

CMD [ "/go-my-fs" ]