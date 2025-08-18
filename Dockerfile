FROM golang:1.24.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd

FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=build /app/auth-service .

EXPOSE 8081

CMD ["./auth-service"]