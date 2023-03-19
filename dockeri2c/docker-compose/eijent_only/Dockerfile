FROM multiarch/ubuntu-core:armhf-bionic As builder
RUN apt-get update &&\
    apt-get install -y curl gcc git
RUN curl -OL https://dl.google.com/go/go1.18.2.linux-armv6l.tar.gz
RUN tar -C /usr/local -xzf go1.18.2.linux-armv6l.tar.gz
RUN rm -rf go1.18.2.linux-armv6l.tar.gz
ENV PATH $PATH:/usr/local/go/bin
# RUN go get -u github.com/davecheney/i2c
# RUN go get -u github.com/d2r2/go-dht
# RUN go get -u github.com/tarm/serial
# RUN go get -u 
WORKDIR /app
ADD ./golang/ ./
RUN go build .
FROM multiarch/ubuntu-core:armhf-bionic
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="0.12b"
LABEL docker_arch="arm"
LABEL name="senser eijent"
LABEL go_version="1.18.2"
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/app /app
CMD ["./app"]
