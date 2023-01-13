FROM golang:1.19.5-alpine AS builder

ADD . /app
WORKDIR /app

RUN go build -o /app/app

FROM alpine:3.17

COPY --from=builder /app/app /app/pipeline

WORKDIR /app

ENTRYPOINT ["./pipeline"]