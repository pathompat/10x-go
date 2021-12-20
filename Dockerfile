# syntax = docker/dockerfile:experimental

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN touch serviceAccount.json
RUN --mount=type=secret,id=auto-devops-build-secrets . /run/secrets/auto-devops-build-secrets && echo $FIRESTORE_SERVICE_ACCOUNT >> serviceAccount.json

COPY *.go ./
COPY ./models ./models

RUN go build -o /docker-10x-go

EXPOSE 5000

CMD [ "/docker-10x-go" ]