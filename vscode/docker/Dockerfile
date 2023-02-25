FROM gitpod/openvscode-server:1.75.1
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="1.75.1"
LABEL go_version="1.19.6"
LABEL name="vscode-server"
USER root
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y curl gcc git locales locales-all make &&\
    locale-gen ja_JP.UTF-8 &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
RUN curl -OL https://dl.google.com/go/go1.19.6.linux-arm64.tar.gz &&\
    tar -C /usr/local -xzf go1.19.6.linux-arm64.tar.gz &&\
    rm -rf go1.19.6.linux-arm64.tar.gz
USER openvscode-server
ENV PATH $PATH:/usr/local/go/bin