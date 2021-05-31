#!/bin/bash

set -e -o pipefail
set -x

# functions
git_setup() {
  cat <<-EOF >"${HOME}"/.netrc
		machine github.com
		login $GITHUB_ACTOR
		password $GITHUB_TOKEN
		machine api.github.com
		login $GITHUB_ACTOR
		password $GITHUB_TOKEN
EOF
  chmod 600 "${HOME}"/.netrc
}
git_set_user() {
  git config --global user.name "$GITHUB_ACTOR"
  git config --global user.email "$GITHUB_EMAIL"
  git config --global committer.name "$GITHUB_ACTOR"
  git config --global committer.email "$GITHUB_EMAIL"
  git config --global author.name "$GITHUB_ACTOR"
  git config --global author.email "$GITHUB_EMAIL"
}

# config gpg
GPG_ASC_FILE="${1:-/init-data/GPG_ASC_FILE}" # /init-data/GPG_ASC_FILE use erda app deploy config; local debug can use $1.
gpg2 -v --batch --import < "${GPG_ASC_FILE}"
git config --global commit.gpgsign true
git config --global user.signingkey "$(gpg --list-secret-keys --keyid-format LONG | grep sec | cut -d' ' -f 4 | cut -d'/' -f 2)"

# config git
GITHUB_ACTOR="${GITHUB_ACTOR:-erda-bot}"
GITHUB_EMAIL="${GITHUB_EMAIL:-erda@terminus.io}"
git_setup
git_set_user

# run bot
/bot