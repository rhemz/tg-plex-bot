FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o plex-bot -ldflags="-w -s" .


FROM scratch
MAINTAINER russ@russ.wtf
ENV GIN_MODE=release
ENV PORT=8080
WORKDIR /app

COPY --from=builder /app/plex-bot .
COPY --from=builder /app/config/config.yaml ./config/config.yaml
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE $PORT
ENTRYPOINT ["./plex-bot"]