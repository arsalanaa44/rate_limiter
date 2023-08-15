FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /rate_limiter

FROM alpine:latest as release


COPY --from=builder /rate_limiter .

EXPOSE 8080

ENTRYPOINT ["./rate_limiter"]
