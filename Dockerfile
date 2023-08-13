FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./


RUN go mod download

COPY  . .

RUN CGO_ENABLED=0 GOOS=windows go build -o /docker-gs-ping


EXPOSE 8080

CMD ["/docker-gs-ping"]

