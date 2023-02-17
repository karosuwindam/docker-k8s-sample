FROM multiarch/ubuntu-core:armhf-bionic As builder
#RUN apt update && apt upgrade -y
RUN apt-get update && \
#    apt-get upgrade -y && \
    apt-get install -y curl gcc git && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
RUN curl -OL https://dl.google.com/go/go1.14.4.linux-armv6l.tar.gz && \
    tar -C /usr/local -xzf go1.14.4.linux-armv6l.tar.gz && \
    rm -rf go1.14.4.linux-armv6l.tar.gz
ENV PATH $PATH:/usr/local/go/bin
RUN go get -u github.com/stianeikeland/go-rpio
WORKDIR /app
ADD ./golang ./
RUN go build -o app
FROM multiarch/ubuntu-core:armhf-bionic
RUN apt-get update && \
#    apt-get upgrade -y && \
    apt-get install -y curl tzdata && \
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
# COPY ./golang/index.html /app
COPY ./golang/conf/ /app/conf/
COPY --from=builder /app/app /app
CMD ["./app"]
