# FROM acentera-cognito-plugin_core
# FROM acentera_core
FROM acentera/prod:golang-base-compile-0.0.2
# registry.gitlab.com/acentera/docker-acentera/serverless/golang-base-build:0.0.1

RUN pip install yq
# RUN apk add --no-cache git

# RUN cd /tmp && \
#     git clone 'https://github.com/stedolan/jq.git' && \
#     sudo -i apt update && \
#     sudo -i apt-get -y install autoconf && \
#     sudo -i apt-get -y install libtool && \
#     sudo -i pip install pipenv && \
#     sudo -i pip install yq && \
#     sudo -i apt-get -y install jq
# RUN apk add --no-cache jq
#  yq

# RUN go get github.com/mitchellh/mapstructure
RUN mkdir -p /go/src/github.com/myplugin/gofaas
RUN mkdir -p /go/src/github.com/acenteracms/acentralib
WORKDIR /go/src/github.com/myplugin/gofaas
