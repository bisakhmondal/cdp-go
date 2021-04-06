FROM golang:1.15-alpine AS build-env
LABEL maintainer="Bisakh Mondal <bisakhmondal00@gmail.com>"

WORKDIR /go-cdp
COPY ./ ./

RUN go build main.go

FROM zenika/alpine-chrome:latest

USER root
RUN apk add bash

WORKDIR /go-cdp
COPY --from=build-env /go-cdp /go-cdp

CMD ["chromium-browser", "--headless", "--no-sandbox", "--remote-debugging-port=9222"]
