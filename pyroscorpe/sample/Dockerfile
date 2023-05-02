FROM golang:1.18.2-buster As builder
WORKDIR /app
ADD ./golang .
RUN go build -o app

FROM debian:unstable-slim
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="0.0.1"
LABEL go_version="1.18.2"
LABEL name="pyroscorpe-sample"

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y locales locales-all curl && \
   locale-gen ja_JP.UTF-8 &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/app /app

CMD ["./app"]