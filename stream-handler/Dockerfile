FROM golang:1.20

WORKDIR /auth

COPY . /auth
COPY go.mod go.sum /auth/
RUN go mod tidy
