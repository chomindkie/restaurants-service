# syntax=docker/dockerfile:1

FROM golang:1.18-alpine

WORKDIR /app
#COPY ./configs ./configs

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /goapp

EXPOSE 8443

CMD [ "/goapp" ]

