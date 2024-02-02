#!/bin/sh
echo "prepare environment"
# replace {% MANAGE_RKEYS_URL %} by content of MANAGE_RKEYS_URL variable
find . -regex '.*\.\(html\|js\|mjs\)$' -print -exec sed -i "s|{% MANAGE_RKEYS_URL %}|${MANAGE_RKEYS_URL}|g" {} \;

echo "start rkeys"

"/srv/rkeys"
