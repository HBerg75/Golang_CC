# syntax=docker/dockerfile:1

FROM golang:1.19 

WORKDIR /app 

COPY go.mod go.sum ./
COPY . .

RUN go build -o /app/main .
CMD ["/app/main"]