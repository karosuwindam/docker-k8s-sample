FROM golang:1.19 As builder

WORKDIR /app
RUN apt-get update && \
    apt-get upgrade -y
ADD ./ .
RUN go build -o app

FROM debian:11
LABEL version="0.0.3"
LABEL name="gotext-server"
LABEL goversion="1.19"
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y locales locales-all curl poppler-utils && \
    locale-gen ja_JP.UTF-8 &&\
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
RUN chown 1000:1000 /app
ENV PYROSCOPE_FLAG false
ADD ./run.sh ./
RUN chmod +x run.sh
USER 1000
RUN mkdir -p public upload/pdf upload/zip db tmp html
ADD ./html ./html
ADD ./txt-tmp ./txt-tmp
COPY --from=builder /app/app /app
CMD ["./run.sh"]