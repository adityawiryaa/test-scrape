FROM golang:1.26-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /controller ./cmd/controller

FROM alpine:3.19

RUN apk add --no-cache ca-certificates
COPY --from=builder /controller /controller

EXPOSE 6001
ENTRYPOINT ["/controller"]
