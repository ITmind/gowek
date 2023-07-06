FROM golang:alpine AS builder
ENV CGO_ENABLED 1
RUN apk update && apk add --no-cache git ca-certificates tzdata upx && update-ca-certificates
RUN apk add --update alpine-sdk
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -ldflags "-s -w -extldflags '-static'" -o app
RUN upx ./app

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo/
COPY --from=builder /build/hds /hds
COPY --from=builder /build/templates /templates
COPY --from=builder /build/static /static
CMD ["/app"]