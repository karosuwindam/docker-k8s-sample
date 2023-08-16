FROM nextcloud:27.0.0
LABEL maintainer "karosu <karosuwindam@gmail.com>"
LABEL version="27.0.0"
LABEL name="nextcloud"
# RUN apk add --no-cache --virtual smbclient 
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install smbclient -y && \
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/entrypoint.sh"]
# CMD ["php-fpm"]
CMD ["apache2-foreground"]
