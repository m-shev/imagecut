FROM golang:1.13-alpine as builder
MAINTAINER mshev
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -ldflags="-w -s" -mod=vendor -o imagecut ./cmd/main.go


FROM alpine:3.10
COPY --from=builder /build/imagecut /opt/imagecut/
COPY --from=builder /build/config/ /opt/imagecut/config
WORKDIR /opt/imagecut
RUN mkdir app-data && mkdir app-data/cache && mkdir app-data/images && mkdir app-data/log
RUN apk update && apk add bash
RUN apk add --no-cache bash
RUN chmod +x ./imagecut
CMD ["./imagecut"]