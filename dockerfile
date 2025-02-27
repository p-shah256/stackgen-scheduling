FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o meeting-scheduler .

FROM alpine:3.17

WORKDIR /app

COPY --from=builder /app/meeting-scheduler .

EXPOSE 8080

CMD ["./meeting-scheduler"]
