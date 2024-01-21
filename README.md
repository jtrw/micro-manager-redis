# micro-manager-caches-keys

Nginx configuration
```
location /index.html {
         proxy_redirect          off;
         proxy_set_header        Host $http_host;
         proxy_pass              http://127.0.1:8080/web/index.html;

```
