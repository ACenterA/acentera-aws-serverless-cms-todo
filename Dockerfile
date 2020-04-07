ARG USER
ARG UID
ARG GID

FROM acentera/prod:node-serverless-cms-base-v0.0.1

ARG USER
ARG UID
ARG GID

RUN /usr/local/bin/userfix.sh && mkdir -p /usr/node_modules && \
    chown -R "${USER}:" /usr/node_modules /root
RUN apk add git
RUN npm i -g cross-env webpack-dev-server@3.7.1
RUN npm i cross-env webpack-dev-server@3.7.1

COPY packag*.json /usr/

RUN chown ${USER}: /usr/ /usr/app
RUN su -m -c "cd /usr/; npm i --prefer-offline --no-audit" - "${USER}"
RUN ls -altrh /usr/node_modules;

# COPY favicon.ico favicon.ico
# COPY .eslintignore .eslintignore
# COPY .eslintrc.js .eslintrc.js
# COPY .babelrc .babelrc
# COPY .postcssrc.js .postcssrc.js
# COPY index-template.html index.html
# 
# COPY build build
# COPY config config
# COPY static static
# COPY src src
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

ENV NODE_PATH=.:/usr/node_modules:/usr/local/node_modules
WORKDIR /usr/app

EXPOSE 9527
EXPOSE 80

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]
CMD ["npm","run", "dev"]
