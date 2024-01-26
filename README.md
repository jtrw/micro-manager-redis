# micro-manager-caches-keys

[![Build](https://github.com/jtrw/micro-manager-redis/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/jtrw/micro-manager-redis/actions)
[![codecov](https://codecov.io/gh/jtrw/micro-manager-redis/graph/badge.svg?token=R4WZPK20B7)](https://codecov.io/gh/jtrw/micro-manager-redis)

Nginx configuration
```
location /index.html {
         proxy_redirect          off;
         proxy_set_header        Host $http_host;
         proxy_pass              http://127.0.1:8080/web/index.html;

```
