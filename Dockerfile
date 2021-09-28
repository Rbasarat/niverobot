# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download

COPY *.go ./

RUN go build -o /niverobot

CMD [ "/niverobot" ]