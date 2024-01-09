#!/bin/bash

# Create PR to Publish analytics-operator to community-operators-prod
#
# * Build and Push operator and bundle images and catalog images
# * Push bundle artifacts to a community-operators-prod branch so that a PR can be manually opened from it.
#
# Usage:
#   VERSION is the tag for the operator image in quay,PREV_VERSION is the version we want to replace
#   FORK is your personal FORK from  community-operators-prod repository 
#   $ VERSION=1.1.2 PREV_VERSION=1.1.1 FORK=<your-fork> PRCREATE=<true/false> ./hack/create_operatorhub_pr.sh
#

set -e

VERSION=${VERSION:-$(make print-VERSION)}
PREV_VERSION=${PREV_VERSION:-$(make print-PREV_VERSION)}
PRCREATE=${PRCREATE:-false}

BRANCH=publish-observability-analytics-operator-$VERSION
FORK=${FORK:-analytics-auto}
GITHUB_TOKEN=${GITHUB_TOKEN:-$ANALYTICS_AUTO_GITHUB_TOKEN}

IMG_REPOSITORY=${IMG_REPOSITORY:-quay.io/openshiftanalytics}

OPERATOR_IMG=$IMG_REPOSITORY/observability-analytics-operator:$VERSION
CATALOG_IMG=$IMG_REPOSITORY/observability-analytics-operator-catalog:$VERSION
BUNDLE_IMG=$IMG_REPOSITORY/observability-analytics-operator-bundle:$VERSION

COMMUNITY_OPERATOR_PROD_GITHUB_ORG=${COMMUNITY_OPERATOR_PROD_GITHUB_ORG:-redhat-openshift-ecosystem}

#If you want to build the operator, and dont want to use the image fro quay uncomment the following line
#make 

# Build bundle directory
make bundle IMG=$OPERATOR_IMG

# Build bundle and catalog images
make bundle-build bundle-push
make catalog-build catalog-push

# Add replaces to dependency graph for upgrade path
if ! grep -qF 'replaces: anomaly-operator.v${PREV_VERSION}' bundle/manifests/anomaly-operator.clusterserviceversion.yaml; then
  sed -i.bak -e "/version: ${VERSION}/a \\
  replaces: anomaly-operator.v$PREV_VERSION" bundle/manifests/anomaly-operator.clusterserviceversion.yaml
fi


# Set Openshift Support Range (bump minKubeVersion in CSV when changing)
if ! grep -qF 'openshift.versions' bundle/metadata/annotations.yaml; then
  sed -i.bak -e "/annotations:/a \\
  com.redhat.openshift.versions: v4.12" bundle/metadata/annotations.yaml
fi

echo "-- Create branch on community-operators-prod fork --"
git clone https://github.com/$COMMUNITY_OPERATOR_PROD_GITHUB_ORG/community-operators-prod.git

mkdir -p community-operators-prod/operators/observability-analytics-operator/$VERSION/
cp -r bundle/* community-operators-prod/operators/observability-analytics-operator/$VERSION/
pushd community-operators-prod/operators/observability-analytics-operator/$VERSION/

git checkout -b $BRANCH
git add ./
git status

message='operator [N] [CI] observability-analytics-operator'
commitMessage="${message} ${VERSION}"
git commit -m "$commitMessage" -s

git remote add upstream git@github.com:$FORK/community-operators-prod.git 

git push upstream --delete $BRANCH || true
git push upstream $BRANCH

if $PRCREATE == "true"; then
  echo "creating pr"
  gh pr create \
    --title "operator observability-analytics-operator (${VERSION})" \
    --body "operator observability-analytics-operator (${VERSION})" \
    --base main \
    --head $FORK:$BRANCH \
    --repo $COMMUNITY_OPERATOR_PROD_GITHUB_ORG/community-operators-prod
fi
popd
