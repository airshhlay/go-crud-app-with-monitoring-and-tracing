# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# final stage
FROM scratch
COPY --from=builder /app/itemService /app/
COPY config ./config
COPY log ./log
COPY --from=builder /etc/ssl/certs ./etc/ssl/certs
# OPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 7000
ENTRYPOINT ["/app/itemService"]