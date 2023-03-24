FROM golang:1.18.2-buster As builder
#RUN apt update && apt upgrade -y
WORKDIR /app
ADD ./golang-v3 .
RUN go build -o app


FROM debian:unstable-slim
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="0.14.6"
LABEL name="booknewread"
# RUN <<EOF
# apt-get update
# apt-get upgrade -y
# apt-get install -y locales locales-all curl
# apt-get clean
# rm -rf /var/lib/apt/lists/*
# EOF
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y locales locales-all curl && \
#    locale-gen ja_JP.UTF-8 &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY ./golang-v3/html html
COPY ./golang-v3/run.sh ./
COPY ./golang-v3/bookmark bookmark
COPY --from=builder /app/app /app
RUN chmod +x run.sh

CMD ["./run.sh"]