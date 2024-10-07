FROM golang:1.22-alpine

WORKDIR /auth

COPY . /auth
COPY go.mod go.sum /auth/
RUN go mod tidy
