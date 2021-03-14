# -------- modules --------
FROM golang:1.16 AS modules

COPY ./go.mod ./go.sum /
RUN go mod download

# -------- build binary --------
FROM golang:1.16 AS builder

RUN useradd -u 1001 appuser

COPY --from=modules /go/pkg /go/pkg
COPY . /build
WORKDIR /build

RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 \
    go build -o ./bin/orderservice ./cmd/orderservice

RUN chmod +x ./bin/orderservice

# -------- build image --------make
FROM scratch


COPY --from=builder /etc/passwd /etc/passwd
USER appuser

COPY --from=builder /build/bin/orderservice /app/bin/orderservice

EXPOSE ${SERVER_PORT}

CMD ["/app/bin/orderservice"]