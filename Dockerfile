FROM golang:1.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.10

RUN adduser -DH hmtpk_schedule

WORKDIR /app

COPY --from=builder /app/main /app/

COPY etc/config.docker.yaml /app/etc/config.yaml
RUN chown hmtpk_schedule:hmtpk_schedule /app
RUN chmod +x /app

USER hmtpk_schedule

CMD ["/app/main"]