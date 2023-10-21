FROM alpine:3.18.3 as builder

WORKDIR /build

RUN apk add ca-certificates cmake git make gcc pkgconfig openssl-dev alpine-sdk &&\
    git clone https://github.com/quicr/quicrq/ &&\
    cd quicrq &&\
    cmake . &&\
    make

FROM alpine
WORKDIR /app
RUN mkdir qlog
COPY --from=builder /build/quicrq/quicrq_app .
