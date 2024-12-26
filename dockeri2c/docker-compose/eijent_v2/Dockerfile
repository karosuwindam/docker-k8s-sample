FROM multiarch/ubuntu-core:arm64-focal As builder
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y curl gcc git locales locales-all &&\
    locale-gen ja_JP.UTF-8
RUN curl -OL https://dl.google.com/go/go1.22.7.linux-arm64.tar.gz &&\
    tar -C /usr/local -xzf go1.22.7.linux-arm64.tar.gz &&\
    rm -rf go1.22.7.linux-arm64.tar.gz
ENV PATH $PATH:/usr/local/go/bin
WORKDIR /app
ADD ./ ./
RUN go build .
FROM multiarch/ubuntu-core:arm64-focal
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="0.14.6"
LABEL name="LABEL_NAME"
LABEL go_version="1.22.7"
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y locales locales-all curl && \
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/app /app
CMD ["./app"]
