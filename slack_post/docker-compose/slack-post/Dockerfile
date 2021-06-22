FROM multiarch/ubuntu-core:armhf-bionic As builder
RUN apt-get update &&\
    apt install -y curl gcc git locales locales-all &&\
    locale-gen ja_JP.UTF-8
RUN curl -OL https://dl.google.com/go/go1.16.5.linux-armv6l.tar.gz &&\
    tar -C /usr/local -xzf go1.16.5.linux-armv6l.tar.gz &&\
    rm -rf go1.16.5.linux-armv6l.tar.gz
ENV PATH $PATH:/usr/local/go/bin
# RUN go get -u 
WORKDIR /app
ADD ./golang/ ./
RUN go build -o app
FROM multiarch/ubuntu-core:armhf-bionic
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y locales locales-all curl && \
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="0.4"
LABEL name="slack-post"
LABEL test="Slackに送信するプログラム<br>webhookを使用してSlackに送信する<br>SLACK_URL:webhook<br>SLACK_CHANNEL:送信チャンネル<br>SLACK_USERNAME:ユーザ名<br>SLACK_COUNTTIME:送信周期"
WORKDIR /app
COPY --from=builder /app/app /app
RUN mkdir tmp
CMD ["./app"]
