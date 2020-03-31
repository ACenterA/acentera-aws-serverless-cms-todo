#!/bin/sh
if [ -n "${SITE_TITLE}" ]; then
  [ -e dist/index.html ] && sed -ri "s/%%TITLE%%/${SITE_TITLE}/g" dist/index.html
  [ -e index-template.html ] && cat index-template.html | sed -r "s/%%TITLE%%/${SITE_TITLE}/g" > index.html
fi
chmod o+rw index.html
[ -n "${USER}" ] && chown ${USER} index.html

[ ! -e index.html ] && cat index-template.html > index.html

# chmod o+rwx node_modules/
# 
# set -x
# [ ! -e 'node_modules/node-module-done' ] && (npm i && touch 'node_modules/node-module-done')
[ -e node_modules ] && rm -fr node_modules
ln -snf ../node_modules
# su -m -c '
# npm i
# ' - ${USER}

# [ ! -e 'node_modules/node-sass/vendor/linux_musl-x64-72/binding.node' ] && npm rebuild node-sass
# #npm rebuild node-sass

export NODE_PATH=${NODE_PATH}:.:/usr/node_modules/:/usr/local/node_modules/

cat > /tmp/.run <<EOF
cd /usr/;
npm ci --prefer-offline --no-audit;
cd app;
$@
EOF
chmod +x /tmp/.run

# rm -fr node_modules;
# cp -Rfp /usr/node_modules .

echo  /bin/sh -m -c \"cd /usr/app; $@\" - \"${USER}\"
exec /bin/sh -m -c "cd /usr/app; ls -latrh; /tmp/.run" - ${USER}
