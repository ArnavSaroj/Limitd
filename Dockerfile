# we define which base image we are going to choose to build our application on

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limiter ./cmd/server

#runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/rate-limiter .

EXPOSE 8080

CMD [ "./rate-limiter" ]
