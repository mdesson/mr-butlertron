### x86 build ###
#FROM golang:latest AS builder
#RUN rm /bin/sh && ln -s /bin/bash /bin/sh
#
#WORKDIR /usr/src/mr-bultertron
#COPY . .
#RUN go mod tidy
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/
#
#FROM alpine:latest
#
#ARG telegram_token
#ARG weather_token
#ARG openai_token
#ENV TELEGRAM_BOT_TOKEN=$telegram_token
#ENV WEATHER_TOKEN=$weather_token
#ENV OPENAI_TOKEN=$openai_token
#
#WORKDIR /root/
#COPY --from=builder /usr/src/mr-bultertron/app /root/
#CMD ["./app"]

### arm build ###
FROM --platform=linux/arm/v7 golang:latest AS builder
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

WORKDIR /usr/src/mr-bultertron
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/

FROM --platform=linux/arm/v7 alpine:latest

ADD https://github.com/golang/go/raw/master/lib/time/zoneinfo.zip /zoneinfo.zip
ENV ZONEINFO /zoneinfo.zip

ARG telegram_token
ARG weather_token
ARG openai_token
ENV TELEGRAM_BOT_TOKEN=$telegram_token
ENV WEATHER_TOKEN=$weather_token
ENV OPENAI_TOKEN=$openai_token

WORKDIR /root/
COPY --from=builder /usr/src/mr-bultertron/app /root/
CMD ["./app"]