#!/bin/bash

set -e -o pipefail
set -x

# git clone erda-bot forked repo to current clean temp dir
git clone "${FORKED_GITHUB_REPO}" .
# add upstream remote to repo where webhook triggers
git remote add upstream "${GITHUB_REPO}"
git remote update && git fetch --all
# checkout CHERRY_PICK_TARGET_BRANCH as base
git checkout --track upstream/"${CHERRY_PICK_TARGET_BRANCH}" -b base-for-auto-cherry-pick
# checkout auto-cherry-pick-pr branch from base
branchForCherryPick="auto-cherry-pick-pr/${GITHUB_PR_NUM}"
git checkout -b "${branchForCherryPick}"
# cherry-pick commit
cherryPickFailedDetailFile="${CHERRY_PICK_FAILED_DETAIL_FILE}"
(git cherry-pick "${MERGE_COMMIT_SHA}" > "${CHERRY_PICK_FAILED_DETAIL_FILE}" | tee) || (git diff >> "${CHERRY_PICK_FAILED_DETAIL_FILE}" && false)
# push to forked repo
git push origin "${branchForCherryPick}" --force
# use hub to create pr
git config -l
# squashed commit
squashedCommitMessage=$(git log -1 --format=medium)
# create pr commit message file
cat <<-EOF >commit_message
Automated cherry pick of #${GITHUB_PR_NUM}: ${PR_TITLE}

Cherry pick of #${GITHUB_PR_NUM} on ${CHERRY_PICK_TARGET_BRANCH}.

Squashed commit message:

\`\`\`
${squashedCommitMessage}
\`\`\`

---

${ORIGIN_ISSUE_BODY}
EOF
# make pull-request
hub pull-request --force -b "erda-project:${CHERRY_PICK_TARGET_BRANCH}" -h "${GITHUB_ACTOR}:${branchForCherryPick}" \
  -l "auto-cherry-pick" --file commit_message > __new_pr_url