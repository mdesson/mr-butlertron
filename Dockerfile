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

FROM golang:latest AS builder
WORKDIR /usr/src/mr-bultertron
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/

FROM alpine:latest
WORKDIR /root/
RUN mkdir ../data/
ENV TELEGRAM_BOT_TOKEN=$TELEGRAM_BOT_TOKEN
ENV WEATHER_TOKEN=$WEATHER_TOKEN
COPY --from=builder /usr/src/mr-bultertron/app /root/
CMD ["./app"]