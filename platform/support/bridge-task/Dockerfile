FROM golang:1.24.7 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o main

FROM alpine
LABEL org.opencontainers.image.source=https://github.com/neopilot-ai/neo
WORKDIR /app
COPY --from=builder /app/main /app/main
CMD ["/app/main"]
