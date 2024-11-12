FROM golang:1.22.1 As builder
#RUN apt update && apt upgrade -y
WORKDIR /app
ADD ./golang .
RUN go build -o app


FROM debian:unstable-slim
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="0.26.1"
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
COPY ./golang/html html
COPY ./golang/run.sh ./
COPY ./golang/bookmark bookmark
COPY --from=builder /app/app /app
ENV PYROSCOPE_FLAG=FALSE
RUN chmod +x run.sh

CMD ["./run.sh"]