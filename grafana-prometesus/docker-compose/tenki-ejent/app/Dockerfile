FROM golang:1.22.1 As builder
WORKDIR /app
ADD ./golang .
RUN go build -o app


FROM debian:unstable-slim
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="0.2.0"
LABEL name="tenki-ejent"
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y locales locales-all curl && \
#    locale-gen ja_JP.UTF-8 &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/app /app

CMD ["./app"]