#!/bin/sh
set -x
export PATH="/usr/local/bin/:${PATH}"
# unset NODE_PATH
export PATH
export NODE_PATH

chmod o+rwx /usr /usr/app
chown ${USER}: /usr /usr/app
cd /usr/app/
[ -e package-lock.json ] && rm package*lock*

if [ -e /usr/app/node_modules/ ]; then
  cp -Rf /usr/local/lib/node_modules/fibers /usr/app/node_modules/
fi

find /usr/app/node_modules -not -user ${USER} -exec chown ${USER} {} \;

cat > /tmp/.run <<EOF
mkdir -p /usr/app/node_modules
cd /usr/app/;
ls -atlrh
whoami
export PATH="/usr/local/bin/:${PATH}"
export PATH="${PATH}"
export NODE_PATH="${NODE_PATH}"
export NODE_ENV="${NODE_ENV}"
npm config set prefix "/usr/app/node_modules"
npm i --no-audit 2>/dev/null;
node /usr/local/proxy.js &
(cd /usr/app/node_modules/fibers/; node-gyp configure; /usr/local/bin/node /usr/app/node_modules/fibers/build)
cd app;
(cd /usr/app; watchexec -f '*package*.json' --force-poll 5000 -i vue-ssr-client-manifest.json,babel.config.json,nodemon.json,package-lock.json --verbose -s SIGKILL -r 'echo Relaunching NPM i; cd /usr/; [ -e package-lock.json ] && rm -f package-lock.json; npm i --prefer-offline --no-audit') &
if [ ! -e dist/ ]; then
  npm run build-all
fi;
$@
EOF
#yarn run build
# yarn run dev
# $@
chmod +x /tmp/.run
npm config set prefix "/usr/app/node_modules"

cd /usr/app;
[ -e node_modules/.cache/babel-loader ] && rm -rf node_modules/.cache/babel-loader

echo "WILL CREATE SYMLINK ln -snf /usr/app/postcss.config.js /usr/postcss.config.js"
export PATH="${PATH}:/usr/app/node_modules/.bin"
export NODE_ENV

exec /tmp/.run
