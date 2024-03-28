FROM golang:1.22.1-alpine AS builder

WORKDIR /usr/src/app
COPY . .
RUN go mod tidy
RUN	go mod download
