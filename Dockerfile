#FROM --platform=linux/arm/v7 golang:latest AS builder
#WORKDIR /usr/src/disc-e
#COPY . .
#RUN go mod tidy
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
#
#FROM --platform=linux/arm/v7 alpine:latest
#FROM alpine:latest
#WORKDIR /root/
#RUN mkdir ../data/
#COPY --from=builder /usr/src/disc-e/app /root/
#COPY --from=builder /usr/src/disc-e/config.json /root/
#CMD ["./app"]

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
ARG opeanai_token
ENV TELEGRAM_BOT_TOKEN=$telegram_token
ENV WEATHER_TOKEN=$weather_token
ENV OPENAI_TOKEN=$opeanai_token

WORKDIR /root/
COPY --from=builder /usr/src/mr-bultertron/app /root/
CMD ["./app"]