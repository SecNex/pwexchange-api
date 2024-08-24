FROM golang:alpine AS builder
WORKDIR /app
COPY src/ /app
RUN go build -o main .

FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app/main

ENV SERVER_SECRET=""
ENV AUTH_TOKEN=""
EXPOSE 8080

CMD ["/app/main"]