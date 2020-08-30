FROM golang
WORKDIR /APP
COPY ./ /APP

FROM resin/rpi-raspbian
RUN apt-get update && apt-get install libraspberrypi-bin -y \
    --no-install-recommends && apt-get clean && rm -rf /var/lib/apt/lists/*