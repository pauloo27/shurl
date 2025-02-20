# STAGE: BUILD
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY Makefile .

COPY internal internal
COPY cmd cmd

RUN go mod download

RUN make dist

# STAGE: TARGET

FROM alpine:latest

RUN addgroup -S user && adduser -S user -G user

USER user

WORKDIR /app
COPY --from=builder /app/shurl /app/shurl

ENTRYPOINT ["/app/shurl"]
