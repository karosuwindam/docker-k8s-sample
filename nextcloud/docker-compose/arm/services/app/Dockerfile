FROM nextcloud
RUN apt-get update&& apt-get upgrade -y
RUN apt-get install smbclient -y
RUN apt-get autoremove -y
ENTRYPOINT ["/entrypoint.sh"]
CMD ["apache2-foreground"]
