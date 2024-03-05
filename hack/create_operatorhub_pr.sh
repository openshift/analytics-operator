#!/bin/bash

# Create PR to Publish analytics-operator to community-operators-prod
#
# * Build and Push operator, bundle and catalog images
# * Push bundle artifacts to a community-operators-prod branch so that a PR can be manually opened from it.
#
# Usage:
#   VERSION is the tag for the operator image in quay,PREV_VERSION is the version we want to replace
#   FORK is your personal FORK from community-operators-prod repository 
#   if you pass it blank then it won't add 'replace' spec in CSV yaml file
#   $ VERSION=1.1.2 PREV_VERSION=1.1.1 FORK=<your-fork> PRCREATE=<true/false> ./hack/create_operatorhub_pr.sh
#

set -e

PROJECT_ROOT="$(git rev-parse --show-toplevel)"
declare -r PROJECT_ROOT

source "$PROJECT_ROOT/hack/utils.bash"
export DOCKER_DEFAULT_PLATFORM=linux/amd64

VERSION=${VERSION:-$(make print-VERSION)}
PREV_VERSION=${PREV_VERSION:-$(make print-PREV_VERSION)}
PRCREATE=${PRCREATE:-false}

BRANCH=publish-observability-analytics-operator-$VERSION
FORK=${FORK:-analytics-auto}
GITHUB_TOKEN=${GITHUB_TOKEN:-$ANALYTICS_AUTO_GITHUB_TOKEN}

COMMUNITY_OPERATOR_PROD_GITHUB_ORG=${COMMUNITY_OPERATOR_PROD_GITHUB_ORG:-redhat-openshift-ecosystem}

info "Remove bundle and community-operators-prod folders"
rm -rf ./bundle
rm -rf ./community-operators-prod

info "Operator image build and push"
make operator-build operator-push VERSION=$VERSION

# Build bundle directory
info "Build Bundle directory"
make bundle VERSION=$VERSION

# Update operator version in CSV
info "Update operator version in CSV"
CSV_FILE_PATH="bundle/manifests/analytics-operator.clusterserviceversion.yaml"
sed -i.bak "s/version: 0.1.0/version: $VERSION/" $CSV_FILE_PATH
sed -i.bak "s/name: analytics-operator.v0.1.0/name: analytics-operator.v$VERSION/" $CSV_FILE_PATH
sed -i.bak "s/observability-analytics-operator:0.1.0/observability-analytics-operator:$VERSION/" $CSV_FILE_PATH

# Add replaces to dependency graph for upgrade path
info "Update 'replaces' for upgrade path"
if ! [[ -z "$PREV_VERSION" ]]; then
  info "condition satisfied"
  sed -i.bak "/version: ${VERSION}/a \\
  replaces: analytics-operator.v$PREV_VERSION" $CSV_FILE_PATH
fi

# Build bundle and catalog images
info "Bundle image build and push"
make bundle-build bundle-push VERSION=$VERSION
info "Catalog image build and push"
make catalog-build catalog-push VERSION=$VERSION

# Set Openshift Support Range (bump minKubeVersion in CSV when changing)
info "Set Openshift Support Range"
ANNOTATIONS_FILE_PATH="bundle/metadata/annotations.yaml"
if ! grep -qF 'openshift.versions' $ANNOTATIONS_FILE_PATH; then
  sed -i.bak -e "/annotations:/a \\
  com.redhat.openshift.versions: v4.14" $ANNOTATIONS_FILE_PATH
fi

info "Remove backup files"
rm -f $CSV_FILE_PATH.bak
rm -f $ANNOTATIONS_FILE_PATH.bak

info "Create branch on community-operators-prod fork"
git clone https://github.com/$COMMUNITY_OPERATOR_PROD_GITHUB_ORG/community-operators-prod.git

mkdir -p community-operators-prod/operators/analytics-operator/$VERSION/
cp -r bundle/* community-operators-prod/operators/analytics-operator/$VERSION/
pushd community-operators-prod/operators/analytics-operator/

git checkout -b $BRANCH
git add ./
git status

message='operator [N] [CI] analytics-operator'
commitMessage="${message} ${VERSION}"
git commit -m "$commitMessage" -s

git remote add upstream git@github.com:$FORK/community-operators-prod.git 

git push upstream --delete $BRANCH || true
git push upstream $BRANCH

if $PRCREATE == "true"; then
  info "creating pr"
  gh pr create \
    --title "operator analytics-operator (${VERSION})" \
    --body "operator analytics-operator (${VERSION})" \
    --base main \
    --head $FORK:$BRANCH \
    --repo $COMMUNITY_OPERATOR_PROD_GITHUB_ORG/community-operators-prod
fi
popd
