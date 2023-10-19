FROM golang:1.21 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /src
COPY . .

RUN go build  .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/azure-janitor /azure-janitor

ENTRYPOINT ["/azure-janitor"]
