FROM acentera/prod:golang-base-compile-0.0.2

RUN pip install yq

RUN go get -u github.com/jinzhu/copier
RUN go get -u github.com/gertd/go-pluralize
RUN go get -u github.com/ACenterA/gqlmerge

RUN mkdir -p /go/src/github.com/myplugin/gofaas
RUN mkdir -p /go/src/github.com/acenteracms/acentralib
WORKDIR /go/src/github.com/myplugin/gofaas
ADD backend/go/Makefile /tmp/Makefile
