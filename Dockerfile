# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY ./models ./models

RUN go build -o /docker-10x-go

EXPOSE 8080

CMD [ "/docker-10x-go" ]