#!/bin/bash

set -e -o pipefail

# already inside temporary dir

# docker login
docker login erda-registry.cn-hangzhou.cr.aliyuncs.com -u "${DOCKER_REGISTRY_USERNAME}" -p "${DOCKER_REGISTRY_PASSWORD}"
# git clone
git clone https://github.com/erda-project/erda-actions -b master
# unshallow
git remote add erda-bot https://github.com/erda-bot/erda-actions.git
git remote update
# create branch
autoBranch="auto-actions-by-issue-comment-${AUTO_BRANCH}"
git checkout -b "${autoBranch}"
# auto build&push image, and git commit
AUTO_GIT_COMMIT=true make "${ACTIONS_TO_MAKE}"
git push erda-bot "${autoBranch}"
# make pr
hub pull-request --force -b "erda-project:master" -l "auto-created-by-pipeline" -h "erda-bot:${autoBranch}" -m "Auto update actions: ${ACTIONS_TO_MAKE}"
