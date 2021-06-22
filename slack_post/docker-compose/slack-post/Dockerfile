FROM multiarch/ubuntu-core:armhf-bionic As builder
RUN apt-get update &&\
    apt-get install -y curl gcc git
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
    apt-get install curl -y && \
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/app /app
RUN mkdir tmp
CMD ["./app"]
