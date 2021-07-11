# Base image
FROM golang:1.15.4-alpine3.12 AS builder
WORKDIR /go/src/course-fetcher
ADD go.* ./
RUN go mod download
ADD . .
RUN mv .config.json.example .config.json
RUN go build

# Starting API
FROM alpine:3.12
WORKDIR /go/src/course-fetcher
COPY --from=builder /go/src/course-fetcher .
CMD ["./course-fetcher", "serve"]