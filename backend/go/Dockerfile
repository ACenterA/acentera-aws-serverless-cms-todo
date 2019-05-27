# FROM acentera-cognito-plugin_core
# FROM acentera_core
FROM acentera/prod:golang-base-compile-0.0.2
# registry.gitlab.com/acentera/docker-acentera/serverless/golang-base-build:0.0.1

RUN pip install yq
# RUN apk add --no-cache git
RUN go get -u github.com/jinzhu/copier
RUN go get -u github.com/gertd/go-pluralize
RUN go get -u github.com/ACenterA/gqlmerge

RUN mkdir -p /go/src/github.com/myplugin/gofaas
RUN mkdir -p /go/src/github.com/acenteracms/acentralib
WORKDIR /go/src/github.com/myplugin/gofaas
