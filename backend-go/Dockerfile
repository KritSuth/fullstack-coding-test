FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG CMD=cmd/api
RUN go build -o server ./${CMD}

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080 50051
CMD ["./server"]
