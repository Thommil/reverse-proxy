FROM nginx:latest

COPY sites-enabled/* /etc/nginx/conf.d/

COPY ssl/server.crt /etc/nginx/conf.d/server.crt
COPY ssl/server.key /etc/nginx/conf.d/server.key

RUN mkdir -p /var/www/html
COPY html/* /var/www/html/