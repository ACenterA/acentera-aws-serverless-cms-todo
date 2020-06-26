export AWS_DEFAULT_REGION ?= us-east-1
APP ?= gofaas

app: dev

clean:
	rm -f $(wildcard handlers/*/main)
	rm -rf $(wildcard web/handlers/*/node_modules)

deploy: BUCKET = pkgs-$(shell aws sts get-caller-identity --output text --query 'Account')-$(AWS_DEFAULT_REGION)
deploy: PARAMS ?= =
deploy: handlers
	@aws s3api head-bucket --bucket $(BUCKET) || aws s3 mb s3://$(BUCKET) --region $(AWS_DEFAULT_REGION)
	sam package --output-template-file out.yml --s3-bucket $(BUCKET) --template-file template.yml
	sam deploy --capabilities CAPABILITY_NAMED_IAM --parameter-overrides $(PARAMS) --template-file out.yml --stack-name $(APP)
	make deploy-static

deploy-static: API_URL=$(shell aws cloudformation describe-stacks --output text --query 'Stacks[].Outputs[?OutputKey==`ApiUrl`].{Value:OutputValue}' --stack-name $(APP))
deploy-static: BUCKET=$(shell aws cloudformation describe-stack-resources --output text --query 'StackResources[?LogicalResourceId==`WebBucket`].{Id:PhysicalResourceId}' --stack-name $(APP))
deploy-static: DIST=$(shell aws cloudformation describe-stack-resources --output text --query 'StackResources[?LogicalResourceId==`WebDistribution`].{Id:PhysicalResourceId}' --stack-name $(APP))
deploy-static: web/static/index.html
	echo "const API_URL=\"$(API_URL)\";" > web/static/js/env.js
	aws s3 sync web/static s3://$(BUCKET)/
	[ -n "$(DIST)" ] && aws cloudfront create-invalidation --distribution-id $(DIST) --paths '/*' || true
	aws cloudformation describe-stacks --output text --query 'Stacks[*].Outputs' --stack-name $(APP)

regen_graphql: build-ui
	([ -e schema.graphql ] && rm -fr schema.graphql) || true
	([ -e sam/graphql.sam ] && cp -f sam/graphql.sam .) || true
	cp -f graphql.sam graphql.tmpl && chmod o+rw graphql.tmpl
	gqlmerge ./schema schema.graphql
	cp -f graphql.sam template.yml && chmod o+rw template.yml
	cp -f graphql.sam /out/template.yml && chmod o+rw /out/template.yml
	# bad cp -f sam/graphql.sam template.yml && chmod o+rw template.yml
	# bad cp -f sam/graphql.sam /out/template.yml && chmod o+rw /out/template.yml
	cp -f schema.graphql /out/.  && chmod o+rw /out/schema.graphql
	cp -f models.sam handlers/app/conf.yaml && chmod o+rw handlers/app/conf.yaml
	echo "FAST ENTRYPOINT IS : ${FAST_ENTRYPOINT}"
	cp -f models.sam ${FAST_ENTRYPOINT}/conf.yaml && chmod o+rw ${FAST_ENTRYPOINT}/conf.yaml
	cp -f template.yml /dist/. || true
	cp -Rfp /go/src/github.com/myplugin/gofaas/handlers /dist/. || true
	cp -f schema.graphql /dist/.

cpy:
	cd /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/ && GOOS=linux go build -ldflags "-s -w" -o main . && cd -
	cp -Rfp /go/src/github.com/myplugin/gofaas/*handlers* /dist/. || true
	find /dist/ -iname '*.go' -exec rm {} \;
	cp -f schema.graphql /dist/. || true
	cp -f template.yml /dist/. || true

dev: regen_graphql handlers cpy
	export GOOS=linux
	mkdir -p /opt; ln -snf /go/src/github.com/acenteracms/acenteralib/aws.so /opt/aws.so
	#make -j dev-watch dev-sam dev-local

dev-go: cpy
	echo "OK DEV-GO CALLED"
	echo "Test RMRTETST2";
	#I=1 && while [ ! /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main -nt /go/src/github.com/myplugin/gofaas/.main ]; do \
	#	echo "Test RMRTETST3"; \
	#	echo "Watiting for compilation..."; \
	#	sleep 5; \
	#	I=$$((I+1)); \
	#	if [ "$${I}" -gt 5 ]; then \
        #         break; \
        #        fi; \
	#done
	#cp -Rfp /go/src/github.com/myplugin/gofaas/*handlers* /dist/. || true
	#find /dist/ -iname '*.go' -exec rm {} \;
	#cp -f schema.graphql /dist/. || true
	#cp -f template.yml /dist/. || true
	kill -9 "$$(cat /tmp/sam3.pid)" || true
	env > /tmp/kkkkkkkl && echo "Launching Fast GO" && { SKIP_COGNITO=true SITE=${PLUGIN_KEY} AWS_SAM_LOCAL=true CUSTOM_PORT=3003 /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main & }; echo $$! > /tmp/sam3.pid
	#env > /tmp/kkkkkkkl && echo "Launching Fast GO" && { SKIP_COGNITO=true SITE=${PLUGIN_KEY} AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY:-'a'} AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID:-'a'} AWS_SAM_LOCAL=true CUSTOM_PORT=3003 /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main & }; echo $$! > /tmp/sam3.pid
	ps auxwww | grep main
	cat /tmp/sam3.pid
	echo 'DONE'

dev-sam:
	kill -9 "$$(cat /tmp/sam.pid)" || true
	echo "Start lambda" && { sam local start-lambda --host 0.0.0.0 -p 3001 --docker-volume-basedir "${HOME_PWD}/" --docker-network ${DOCKER_NETWORK} & echo $$! > /tmp/sam1.pid; }
	kill -9 "$$(cat /tmp/sam2.pid)" || true
	echo "Launching INVOKE" && { sam local start-api --force-image-build --debug --host 0.0.0.0 -p 3000 --docker-volume-basedir "${HOME_PWD}/" --docker-network ${DOCKER_NETWORK} -s web/static & echo $$! > /tmp/sam2.pid; }
	#kill -9 "$$(cat /tmp/sam3.pid)" || true
	#echo "Launching Fast GO" && { CUSTOM_PORT=3003 /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main & echo $$! > /tmp/sam3.pid; }
	#echo "Launching Fast GO" && { CUSTOM_PORT=3003 /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main; }

removeapp:
	if [ -e /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main ]; then
	  rm -f /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main; \
	fi

	# touch /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main.go; \
		#
build-ui:
	(mkdir -p /usr/local/go/src/github.com/acenterastatic/static/ && cd /ui && go-bindata -o /usr/local/go/src/github.com/acenterastatic/static/main.go -pkg static dist/...) || true
	cp -f /usr/local/go/src/github.com/acenterastatic/static/main.go /go/src/github.com/acenterastatic/static/. || true

dev-watch:
	cp -Rfp /shared /dist/. || true
	mkdir -p /opt; ln -snf /go/src/github.com/acenteracms/acenteralib/aws.so /opt/aws.so
	watchexec -f '*.yml' -s SIGKILL -r 'make -j dev-sam' &
	go build -i
	watchexec -f '*.go' -s SIGKILL -r 'make -j handlers; make -j dev-go' &
	find . -maxdepth 4 -type f -iname main.go -exec touch {} \;
	sleep 1; touch ${FAST_ENTRYPOINT}/main.go
	watchexec -w schema -e graphql -r 'make regen_graphql'

#build-all:
#	for f in $(ls -d *handlers*); do find $f -type d -exec (cd ${d};done;

HANDLERS=$(addsuffix main,$(wildcard *handlers*/*/))
$(HANDLERS): handlers/%/main: *.go handlers/%/main.go
#$(HANDLERS): handlers/%/main: *.go handlers/%/main.go vendor
	(ps -ef | grep "${FAST_ENTRYPOINT}/" |grep -v grep | awk '{ print $2 }' |xargs kill ) || true
	# (cd ./$(dir $@) && mv -f main main.$$; rm -f main.$$; cd -) || true
	cd ./$(dir $@) && GOOS=linux go build -ldflags "-s -w" -o main . && cd -
	cd /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/ && GOOS=linux go build -ldflags "-s -w" -o main . && cd -
	if [ "$(dir $@)" = "handlers/website-public/" -o "$(dir $@)" = "handlers/migration/" ]; then cd ./$(dir $@) && GOOS=linux go build -ldflags "-s -w" -tags=test -o main .; else cd ./$(dir $@) && GOOS=linux go build -ldflags "-s -w" -o main .; fi;
	if [ "$(dir $@)" = "handlers/shell/" ]; then cd ./$(dir $@) && env GOOS=linux GOARCH=arm64 GOARM=5 go build -o main-arm .; fi;
	if [ "$(dir $@)" = "handlers/app/" ]; then cd ./$(dir $@) && cp -f *.yaml /go/src/github.com/acenteracms/acenteralib/. && cp -f main /go/src/github.com/acenteracms/acenteralib/.; fi;

HANDLERS_JS=$(addsuffix node_modules,$(wildcard web/handlers/*/))
$(HANDLERS_JS): web/handlers/%/node_modules: web/handlers/%/package.json
	cd ./$(dir $@) && npm install && node-prune >/dev/null && touch node_modules

rmtest:
	if [ -e /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main ]; then \
		echo > /go/src/github.com/myplugin/gofaas/.main; \
	fi
	echo "AAA";
	if [ -e /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main ]; then \
		echo "CCC"; \
		if [ /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main.go -nt /go/src/github.com/myplugin/gofaas/.main ]; then \
			echo "DDD"; \
			touch /go/src/github.com/myplugin/gofaas/${FAST_ENTRYPOINT}/main.go; \
			touch /go/src/github.com/myplugin/gofaas/.main; \
		fi; \
	fi

handlers: rmtest handlers-go handlers-js build-all
handlers-go: $(HANDLERS)
handlers-js: $(HANDLERS_JS)

template: handlers
	cd /go/src/github.com/myplugin/gofaas/ && bash sam-package.sh prod

test: vendor
	go test -v ./...

vendor:
	go build
	#go get github.com/golang/dep/cmd/dep
	#dep ensure
