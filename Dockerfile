FROM golang:alpine as builder

WORKDIR /app
COPY go.mod go.sum config.yaml ./
COPY cmd/ cmd/
COPY internal/ internal/

RUN go mod download &&\
    go build -o server cmd/server/server.go

######
# Server
FROM scratch

WORKDIR /app

COPY --from=builder /app/server ./
COPY --from=builder /app/config.yaml ./

EXPOSE 8080

CMD [ "./server" ]