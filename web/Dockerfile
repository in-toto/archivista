ARG BASE_IMAGE=nginx:1.23.2-alpine
FROM ${BASE_IMAGE}
COPY ./nginx.prod.conf /etc/nginx/conf.d/default.conf
COPY ./public /usr/share/nginx/html