#!/bin/sh

set -e

# if test $# -ne 1
# then
#   echo "Usage: $0 <chart>;<oldver>;<newver>"
#   exit 1
# fi
#
OPTSTRING="f:m:"

while getopts ${OPTSTRING} opt; do
  case ${opt} in
    f)
      # I really need to rewrite this.
      arg=$(just check-versions --filter ${OPTARG} | head -1)
      ;;
    m)
      arg=${OPTARG}
      ;;
    ?)
      echo "Invalid option: -${OPTARG}."
      echo
      echo "Usage: $0 [-m 'prometheus;0.1.0;0.1.2'|-f prom]"
      exit 1
      ;;
  esac
done

if test "${arg}" = ""
then
    echo "Nothing to update."
    exit 1
fi

chart=$(echo $arg | cut -d ';' -f1 | sed -e 's/\//\\\//g')
target_update=$(echo $chart | tr -d '\\')
oldver=$(echo $arg | cut -d ';' -f2)
newver=$(echo $arg | cut -d ';' -f3)

git fetch
git remote prune origin

git checkout versions.yaml

sed -e "s/  $chart: $oldver/  $chart: $newver/" -i versions.yaml

current_branch=$(git rev-parse --abbrev-ref HEAD)
pr_branch_name=$(echo $chart:$oldver | tr -d '\\' | tr '/:' '_')

is_branch_existing=$(git ls-remote --heads origin refs/heads/${pr_branch_name} | wc -l)

if test ${is_branch_existing} -ne 0
then
  echo "remote branch ${is_branch_existing} is already existing."
fi

git diff versions.yaml

git checkout -b ${pr_branch_name}

git commit -m "Updating ${target_update} from $oldver to $newver" versions.yaml
git push --force --set-upstream origin ${pr_branch_name}


if test ${is_branch_existing} -eq 0
then
  tea pr --repo mycroft/k8s-home create \
    --head ${pr_branch_name} \
    --title "${target_update} update to ${newver}" \
    --description "Update from ${oldver} to ${newver}"
else
  echo "skip creating PR as branch was already existing."

  pr_id=$(tea pr --repo mycroft/k8s-home --state open --fields index,head --output simple|grep ${pr_branch_name}|cut -d' ' -f1)

  echo "pr id is ${pr_id}"

  tea comment --repo mycroft/k8s-home ${pr_id} "Branch forced-pushed: Update from ${oldver} to ${newver}"
fi

git checkout ${current_branch}
git branch -D ${pr_branch_name}
git checkout versions.yaml
