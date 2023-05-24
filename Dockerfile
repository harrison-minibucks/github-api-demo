FROM golang:1.20 AS builder

COPY . /src
WORKDIR /src

RUN go env -w GO111MODULE=on
RUN make build

FROM alpine
COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

CMD ["./github-api-demo", "-conf", "/data/conf"]
