# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o acbr-cep-api

# Runtime stage
FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=builder /app/acbr-cep-api .
COPY ./lib/libacbrcep64.so ./lib/

RUN chmod +x ./lib/libacbrcep64.so

EXPOSE 8080

CMD ["./acbr-cep-api"]