#!/bin/bash
# set -e
# set -x

PROGNAME=$(basename -- "${0}")
PROJROOT=$(d=$(dirname -- "${0}"); cd "${d}/.." && pwd)

if [[ -z "$1" ]]; then
  echo "Next 1st parameter dev, qa, or prod for stage"
else

   STAGE=$1
   DT=$(date +%Y)
   SEMVER=0.0.18
   S3PREFIX="packaged/$DT/0.0.1/acentera-${PLUGINNAME}"
   BUCKETNAME=${S3_BUCKET:-"lambda-at-edge-dev-serverlessdeploymentbucket-1gmbbmp4ajnba"}

                                                                                
  if [ -e .${STAGE}.aws ]; then
    source .${STAGE}.aws
  fi

  cp -f template.yml .template.yml.$1                                           
  # First update the Path: .... no hard-coded value ideally                     
  sed -ri "s~<%PLUGIN_NAME%>~${PLUGINNAME}~g" .template.yml.$1                  
  sed -ri "s~<%STAGE%>~${STAGE}~g" .template.yml.$1                  
  sed -ri "s~<%SEMVER%>~${SEMVER}~g" .template.yml.$1                  
  [[ -e packaged-template.yml.$1 ]] && rm -f packaged-template.yml.$1      

  mkdir -p /go/src/github.com/myplugin/gofaas/shared/
  /bin/cp -f /go/src/github.com/acenteracms/acenteralib/aws.so /go/src/github.com/myplugin/gofaas/shared/.
  echo "TEST INVOKE"
  sam local invoke --template .template.yml.$1 "ModelLambda" --docker-volume-basedir "${HOME_PWD}/" -e event.json
  sam package --debug --template-file .template.yml.$1 --output-template-file packaged-template-acentera.yaml.${STAGE} --s3-bucket ${BUCKETNAME} --s3-prefix ${S3PREFIX}

  echo  "############################################################"
  echo  "################# CMS ##############################"
  echo  "############################################################"

  cp -f template.plugin.yml .template.plugin.yml.$1
  sed -ri "s~<%STAGE%>~${STAGE}~g" .template.plugin.yml.$1
  sed -ri "s~${PLUGINNAME}~<%PLUGIN_NAME%>~g" .template.plugin.yml.$1
  sed -ri "s~<%PLUGIN_NAME%>~${PLUGINNAME}~g" .template.plugin.yml.$1
  sed -ri "s~<%SEMVER%>~${SEMVER}~g" .template.plugin.yml.$1

  cat packaged-template-acentera.yaml.${STAGE}

  NEW_LICENCEFILE=$(cat packaged-template-acentera.yaml.prod | yq . | jq -r '.Metadata["AWS::ServerlessRepo::Application"].LicenseUrl' | sed -r 's~(s3://.*/packaged/.*)~\1~g')
  NEW_READMEFILE=$(cat packaged-template-acentera.yaml.prod | yq . | jq -r '.Metadata["AWS::ServerlessRepo::Application"].ReadmeUrl' | sed -r 's~(s3://.*/packaged/.*)~\1~g')

  NEW_S3FILE=$(cat packaged-template-acentera.yaml.${STAGE} | yq . | jq -r '.Resources.RequestsLayer.Properties.ContentUri')
  # NEW_S3FILE=$(cat packaged-template-acentera.yaml.${STAGE} | yq . | jq -r '.Resources.RequestsLayer.Properties.ContentUri' | sed -r 's~s3://(.*)/(packaged/.*)~\2~g')
  echo $NEW_S3FILE
  # cat .template.plugin.yml.$1 | yq . | jq --arg LayerBIN $NEW_S3FILE '.Resources.RequestsLayerCMS.Properties.Content.S3Key = $LayerBIN' > .template.plugin.yml.$1.tmp
  cat .template.plugin.yml.$1 | yq . | jq --arg LayerBIN $NEW_S3FILE '.Resources.RequestsLayerCMS.Properties.ContentUri = $LayerBIN' > .template.plugin.yml.$1.tmp
  cp -f .template.plugin.yml.$1.tmp .template.plugin.yml.$1
  
  SCHEMA_CONTENT=$(cat schema.graphql)
  cat .template.plugin.yml.$1 | yq . | jq --arg LayerBIN "$SCHEMA_CONTENT" '.Resources.AppSyncSchema.Properties.Definition = $LayerBIN' > .template.plugin.yml.$1.tmp
  cp -f .template.plugin.yml.$1.tmp .template.plugin.yml.$1

  # LicenceFile
  cat .template.plugin.yml.$1 | yq . | jq --arg LayerBIN $NEW_LICENCEFILE '.Metadata["AWS::ServerlessRepo::Application"].LicenseUrl = $LayerBIN' > .template.plugin.yml.$1.tmp
  cp -f .template.plugin.yml.$1.tmp .template.plugin.yml.$1
  # ReadmeFile
  cat .template.plugin.yml.$1 | yq . | jq --arg LayerBIN $NEW_READMEFILE '.Metadata["AWS::ServerlessRepo::Application"].ReadmeUrl = $LayerBIN' > .template.plugin.yml.$1.tmp
  cp -f .template.plugin.yml.$1.tmp .template.plugin.yml.$1

  # NEW_S3BUCKET=$(cat packaged-template-acentera.yaml.${STAGE} | yq . | jq -r '.Resources.RequestsLayer.Properties.ContentUri' | sed -r 's~s3://(.*)/(packaged/.*)~\1~g')
  # echo $NEW_S3BUCKET
  # cat .template.plugin.yml.$1 | yq . | jq --arg LayerBIN $NEW_S3BUCKET '.Resources.RequestsLayerCMS.Properties.Content.S3Bucket = $LayerBIN' > .template.plugin.yml.$1.tmp

  cp -f .template.plugin.yml.$1.tmp .template.plugin.yml.$1
  NEW_AS=$(cat packaged-template-acentera.yaml.${STAGE} | yq . | jq -r '.Resources.ApiApp.Properties.CodeUri')
  echo $NEW_AS
  cat .template.plugin.yml.$1 | yq . | jq --arg LayerBIN $NEW_AS '.Resources.ApiApp.Properties.CodeUri = $LayerBIN' > .template.plugin.yml.$1.tmp
  cp -f .template.plugin.yml.$1.tmp .template.plugin.yml.$1

  NEW_AS=$(cat packaged-template-acentera.yaml.${STAGE} | yq . | jq -r '.Resources.ApiPluginSettings.Properties.CodeUri')
  echo $NEW_AS
  cat .template.plugin.yml.$1 | yq . | jq --arg LayerBIN $NEW_AS '.Resources.ApiPluginSettings.Properties.CodeUri = $LayerBIN' > .template.plugin.yml.$1.tmp
  cp -f .template.plugin.yml.$1.tmp .template.plugin.yml.$1

  NEW_AS=$(cat packaged-template-acentera.yaml.${STAGE} | yq . | jq -r '.Resources.ModelLambda.Properties.CodeUri')
  echo $NEW_AS
  cat .template.plugin.yml.$1 | yq . | jq --arg LayerBIN $NEW_AS '.Resources.ModelLambda.Properties.CodeUri = $LayerBIN' > .template.plugin.yml.$1.tmp
  cp -f .template.plugin.yml.$1.tmp .template.plugin.yml.$1

  NEW_AS=$(cat packaged-template-acentera.yaml.${STAGE} | yq . | jq -r '.Resources.PublicWebsite.Properties.CodeUri')
  echo $NEW_AS
  cat .template.plugin.yml.$1 | yq . | jq --arg LayerBIN $NEW_AS '.Resources.PublicWebsite.Properties.CodeUri = $LayerBIN' > .template.plugin.yml.$1.tmp
  cp -f .template.plugin.yml.$1.tmp .template.plugin.yml.$1

  cp -f .template.plugin.yml.$1 output.yml
  # sam publish --debug --template packaged-template-acentera-ecseks-resources.yaml.$1 --region us-east-1
fi
