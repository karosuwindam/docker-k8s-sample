FROM golang:1.16.5-buster As builder
#RUN apt update && apt upgrade -y
WORKDIR /app
ADD ./golang .
RUN go build -o app


FROM debian:unstable-slim
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="0.10.3"
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
COPY ./golang/config config
COPY ./golang/html html
COPY ./golang/bookmarks.html ./
COPY ./golang/run.sh ./
COPY ./golang/bookmark bookmark
COPY --from=builder /app/app /app
RUN chmod +x run.sh

CMD ["./run.sh"]