#!/bin/sh
set -x
if [ ! -e node_modules/webpack ]; then
   ln -snf /var/app/node_modules /usr/app/node_modules
fi
npm rebuild node-sass
exec $@
