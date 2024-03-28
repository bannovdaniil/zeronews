FROM golang:1.22.1-alpine AS builder

WORKDIR /usr/local/src

COPY . .

RUN go mod download
RUN mkdir "./bin"
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -v -o ./bin/zeronews-server ./cmd/zeronews

ENTRYPOINT ["./bin/zeronews-server"]
