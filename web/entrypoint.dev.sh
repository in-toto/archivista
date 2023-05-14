#!/bin/sh
/usr/sbin/nginx -c /etc/nginx/nginx.conf -g "daemon off;" &

npm i
npm run start