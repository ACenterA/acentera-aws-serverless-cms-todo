ARG USER
ARG UID
ARG GID

FROM acentera/prod:acentera-node-python-0.0.3

ARG USER
ARG UID
ARG GID

RUN apk add git make gcc g++

# ENV DEBUG=true
RUN /usr/local/bin/userfix.sh && echo 'ok' && mkdir -p /usr/app /usr/app/node_modules /root/.cache/yarn && \
    chown -R ${USER}: /root/.npm /root/.config /root/.cache/yarn && \
    chown ${USER}: /root /usr/app/node_modules /usr && \
    chown -R ${USER}: /usr/app/node_modules

RUN npm config set prefix "/usr/app/node_modules" && \
    ln -snf /root/.config /home/${USER}/.config && \
    chown -R $USER:$(id -gn $USER) /root/.config && \
    ls -altrh /home/${USER}/

RUN chmod o+rwx /usr /usr/app && chown ${USER}: /usr
RUN chmod o+rwx /usr/app/node_modules && chown ${USER}: /usr/app/node_modules
RUN chown ${USER}: /usr/ /usr/app

VOLUME /usr/app/node_modules
VOLUME /usr/app

ENV PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:.:/usr/local/bin:/bin:/usr/app/node_modules/.bin/:/usr/local/lib/node_modules/npm/bin/:/usr/local/lib/node_modules/npm/bin/node-gyp-bin/
ENV NODE_PATH=ENV NODE_PATH=.:/usr/app/node_modules:/opt/node_modules/:/opt/:/opt/shared/:/opt/node/node_modules/:/usr/local/lib/node_modules:/usr/local/node_modules/lib/node_modules/:/usr/local/node_modules/:/var/task/
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

WORKDIR /usr/app

EXPOSE 9527
EXPOSE 80

ENTRYPOINT [ "/usr/local/bin/entrypoint.sh" ]
CMD ["npm","run", "dev"]
