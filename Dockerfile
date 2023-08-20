FROM golang:alpine

WORKDIR /app
COPY . /app/

RUN go mod download
RUN go build -o server cmd/server/server.go

EXPOSE 8080

CMD [ "./server" ]