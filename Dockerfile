FROM golang:alpine as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o plex-bot -ldflags="-w -s" .


FROM scratch
MAINTAINER russ@russ.wtf
ENV GIN_MODE=release
ENV PORT=8080
WORKDIR /app

COPY --from=builder /app/plex-bot .
EXPOSE $PORT
ENTRYPOINT ["./plex-bot"]