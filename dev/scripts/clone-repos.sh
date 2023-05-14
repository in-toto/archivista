#! /bin/bash
set -ex

GITLABREPOS=(
	"testifysec/judge-platform/scribe" \
	"testifysec/judge-platform/spire-internal" \
	"testifysec/judge-platform/web" \
	"testifysec/judge-platform/judge-api"\
)

GITHUBREPOS=(
	"testifysec/go-witness" \
	"testifysec/witness" \
	"testifysec/archivista" \
	"testifysec/witness-install-action" \
)

function cloneall {
	local REMOTE=$1
	shift
	local REPOS=("$@")
	for repo in "${REPOS[@]}"; do
			##basename, remove git extension
			repo_name=$(basename $repo | sed 's/\.git//g')
			echo "Cloning from $repo to $REPOS_DIR/$repo_name"

			## check to see if repo exists
			if [ -d "$REPOS_DIR/$repo_name" ]; then
					echo "Repo $repo already exists, skipping"
					continue
			fi

			local branchopt=""
			if [[ ! -z "${BRANCH}" ]]; then
				branchopt="--branch ${BRANCH}"
			fi

			git clone "$REMOTE$repo" $REPOS_DIR/$repo_name ${branchopt}
	done
}

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPOS_DIR=$(cd "$SCRIPT_DIR/../../" && pwd)
GITLAB="git@gitlab.com:"
GITHUB="git@github.com:"

while getopts 'h' OPT
do
	case "${OPT}" in
		h) GITHUB="https://github.com/"
			 GITLAB="https://gitlab.com/"
			;;
	esac
done
shift $((OPTIND - 1))

# use first argument as the branch name
BRANCH=$1
cloneall $GITLAB "${GITLABREPOS[@]}"
cloneall $GITHUB "${GITHUBREPOS[@]}"
