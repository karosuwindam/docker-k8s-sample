FROM ubuntu:22.04 As builder
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y curl gcc git locales locales-all &&\
    locale-gen ja_JP.UTF-8
RUN curl -OL https://go.dev/dl/go1.22.6.linux-arm64.tar.gz &&\
    tar -C /usr/local -xzf go1.22.6.linux-arm64.tar.gz &&\
    rm -rf go1.22.6.linux-arm64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
WORKDIR /app
ADD ./golang .
RUN go build -o app
FROM ubuntu:22.04
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="0.5.0"
LABEL docker_arch="arm64"
LABEL go_version="1.22.6"
LABEL name="client-go"

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y locales locales-all curl && \
#    locale-gen ja_JP.UTF-8 &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY ./golang/html html
COPY --from=builder /app/app /app
CMD ["./app"]