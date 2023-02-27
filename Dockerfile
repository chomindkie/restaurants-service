# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /goapp

CMD [ "/goapp" ]

