#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Function to get the latest golangci-lint version for a given release track
get_latest_golangci_lint_version() {
  local releaseTrack=$1
  local url="https://api.github.com/repos/golangci/golangci-lint/releases"
  local releases=$(curl -s $url | jq -r ".[].tag_name")

  for tag in $releases; do
    if [[ $tag == v$releaseTrack.* ]]; then
      echo ${tag#v} # Remove the 'v' prefix
      return 0
    fi
  done

  echo "No matching release found for track: $releaseTrack" >&2
  return 1
}

# Function to update prbuild.yaml
update_prbuild_yaml() {
  local newVersion=$1
  local yamlFile=".github/workflows/prbuild.yaml"
  sed -i.bak -E "/- name: golangci-lint/,/version:/s/(version: )v[0-9]+\.[0-9]+\.[0-9]+/\1v$newVersion/" "$yamlFile"
  rm "${yamlFile}.bak"
  echo "Updated golangci-lint version field in $yamlFile to v$newVersion"
}

# Function to update golangci-lint script
update_golangci_lint_script() {
  local newVersion=$1
  local scriptFile="hack/golangci-lint"
  sed -i.bak -E "s|(golangci-lint@)v[0-9]+\.[0-9]+\.[0-9]+|\1v$newVersion|g" "$scriptFile"
  rm "${scriptFile}.bak"
  echo "Updated $scriptFile to golangci-lint version $newVersion"
}

# Main function
main() {
  if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <releaseTrack>"
    exit 1
  fi

  local releaseTrack=$1
  local latestVersion

  latestVersion=$(get_latest_golangci_lint_version "$releaseTrack")
  if [[ $? -ne 0 ]]; then
    echo "Failed to fetch the latest golangci-lint version for release track $releaseTrack" >&2
    exit 1
  fi

  echo "Latest golangci-lint version for release track $releaseTrack: $latestVersion"

  update_prbuild_yaml "$latestVersion"
  update_golangci_lint_script "$latestVersion"

  echo "golangci-lint version updated successfully to $latestVersion"
}

main "$@"
