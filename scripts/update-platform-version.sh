#!/bin/bash

# update-platform-version
# This script helps us keep our platform version up to date

version_output=$(npx git-conventional-commits version)
latest_version=$(echo "$version_output" | tail -n 1)
current_version=$(node -pe "require('./package.json').version")

# Update the package.json file with the new version
if [[ $latest_version != $current_version ]]; then
  node -e "
  const fs = require('fs');
  const packageJson = JSON.parse(fs.readFileSync('package.json'));
  const packageLockJson = JSON.parse(fs.readFileSync('package-lock.json'));

  packageJson.version = '${latest_version}';
  packageLockJson.version = '${latest_version}';

  fs.writeFileSync('package.json', JSON.stringify(packageJson, null, 2));
  fs.writeFileSync('package-lock.json', JSON.stringify(packageLockJson, null, 2));
  "

  # Stage the updated package.json file
  git add package.json
  git add package-lock.json

  echo "Updated package.json version to ${latest_version} and staged the file."
else
  echo "No version update found in git-conventional-commits."
fi
