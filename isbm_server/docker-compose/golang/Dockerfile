FROM golang AS builder
RUN go get -u github.com/seihmd/openbd
RUN go get -u github.com/PuerkitoBio/goquery
RUN go get -u github.com/go-sql-driver/mysql
RUN go get -u github.com/mattn/go-sqlite3
RUN go get -u github.com/go-delve/delve/cmd/dlv
WORKDIR /app
COPY . .
RUN go build -o app

FROM ubuntu:20.10
WORKDIR /app
COPY ./config config
COPY ./html html
COPY ./tmp_html tmp_html
COPY --from=builder /app/app ./