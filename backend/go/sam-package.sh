#!/bin/bash
set -e
set -x
STAGE=${1-"prod"}
DT=$(date +%Y)
SEMVER=${SEMVER:-$CODEBUILD_SOURCE_VERSION}
S3PREFIX="packaged/$DT/0.0.1/acentera-${PLUGINNAME}"
BUCKETNAME=${S3_BUCKET:-"not-specified"}

if [ -e .${STAGE}.aws ]; then
  source .${STAGE}.aws
fi

# Replace some token
cp -f template.yml .template.yml.$1
# First update the Path: .... no hard-coded value ideally
sed -ri "s~<%PLUGIN_NAME%>~${PLUGINNAME}~g" .template.yml.$1
sed -ri "s~<%STAGE%>~${STAGE}~g" .template.yml.$1
sed -ri "s~<%SEMVER%>~${SEMVER}~g" .template.yml.$1
[[ -e packaged-template.yml.$1 ]] && rm -f packaged-template.yml.$1

mkdir -p /go/src/github.com/myplugin/gofaas/shared/
/bin/cp -f /go/src/github.com/acenteracms/acenteralib/aws.so /go/src/github.com/myplugin/gofaas/shared/.

# Generate lambda layers
sam local invoke --template .template.yml.$1 "ModelLambda" --docker-volume-basedir "${HOME_PWD}/" -e event.json
sam package --debug --template-file .template.yml.$1 --output-template-file output.yml --s3-bucket ${BUCKETNAME} --s3-prefix ${S3PREFIX}
