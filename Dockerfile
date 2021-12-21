# syntax = docker/dockerfile:experimental

FROM golang:1.16-alpine

WORKDIR /app

# set default env
ENV GOPROXY https://goproxy.io
ENV GIN_MODE release

# default port of auto devops is 5000
ENV PORT 5000
EXPOSE 5000

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# initial firestore credential
RUN touch serviceAccount.json
RUN --mount=type=secret,id=auto-devops-build-secrets . /run/secrets/auto-devops-build-secrets && echo $FIRESTORE_SERVICE_ACCOUNT >> serviceAccount.json

COPY *.go ./
COPY ./models ./models

RUN go build -o /docker-10x-go

CMD [ "/docker-10x-go" ]