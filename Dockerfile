FROM debian:bullseye-slim

ENV TZ=Europe/Berlin
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV LANG en_US.utf8
RUN echo "deb http://ftp.debian.org/debian bullseye-backports main" > /etc/apt/sources.list.d/bullseye_backports.list && \
	DEBIAN_FRONTEND=noninteractive apt-get update && \
    apt-get install -y --no-install-recommends apt-utils curl gnupg apt-transport-https supervisor ca-certificates dirmngr locales bzip2 && \
    apt-get install -y --no-install-recommends certbot python3-certbot-nginx -t bullseye-backports && \
	localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8 && \
	apt-get clean && \
	apt-get autoremove -y && \
	rm -rf /var/lib/apt/lists/*

RUN echo "deb http://nginx.org/packages/mainline/debian/ bullseye nginx" > /etc/apt/sources.list.d/nginx.list && \
	curl -sf https://nginx.org/keys/nginx_signing.key | apt-key add - && \
	DEBIAN_FRONTEND=noninteractive apt-get update && \
	apt-get remove -y nginx-common && \
	apt-get install -y --no-install-recommends nginx && \
	apt-get clean && \
	apt-get autoremove -y && \
	rm -rf /var/lib/apt/lists/*

ADD assets /
RUN rm /etc/nginx/sites-enabled/default && \
    mkdir -p /var/lib/letsencrypt/www/ /etc/nginx/generated.d/

EXPOSE 80 443

ADD dist/container-watch /usr/local/bin/container-watch

CMD [ "/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf" ]
