FROM multiarch/ubuntu-core:armhf-bionic As builder
#RUN apt update && apt upgrade -y
RUN apt update
RUN apt install -y curl gcc git
RUN  curl -OL https://dl.google.com/go/go1.14.4.linux-armv6l.tar.gz
RUN tar -C /usr/local -xzf go1.14.4.linux-armv6l.tar.gz
RUN rm -rf go1.14.4.linux-armv6l.tar.gz
ENV PATH $PATH:/usr/local/go/bin
RUN go get -u github.com/jinzhu/gorm
RUN go get -u github.com/davecheney/i2c
WORKDIR /app
ADD ./golang/ ./
RUN go build -o app
FROM multiarch/ubuntu-core:armhf-bionic
WORKDIR /app
COPY --from=builder /app/app /app
CMD ["./app"]
