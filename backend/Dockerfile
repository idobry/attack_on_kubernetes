# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.12 as certs
RUN apk --update add ca-certificates

# final stage
FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/attack_on_kubernetes /app/
EXPOSE 8080
ENTRYPOINT ["/app/attack_on_kubernetes", "start"]