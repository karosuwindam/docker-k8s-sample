FROM nextcloud:23.0.3-fpm-alpine
#RUN apt-get update&& apt-get upgrade -y
#RUN apt-get install smbclient -y
#RUN apt-get autoremove -y
RUN apk add --no-cache --virtual smbclient 
ENTRYPOINT ["/entrypoint.sh"]
CMD ["apache2-foreground"]
