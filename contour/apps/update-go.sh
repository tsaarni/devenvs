#! /usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Define the files to patch with the Go version.
filesToPatchWithGolangVersion=(
  "Makefile"
  ".github/workflows/build_daily.yaml"
  ".github/workflows/build_tag.yaml"
  ".github/workflows/codeql-analysis.yml"
  ".github/workflows/prbuild.yaml"
)

# Function to get the latest Go version.
getLatestGoVersion() {
  local releaseTrack=$1
  local url="https://go.dev/dl/?mode=json&include=all"
  local latestVersion=$(curl -s $url | jq -r --arg track "go$releaseTrack" '.[] | select(.version | startswith($track)) | .version' | head -n 1)

  if [ -z "$latestVersion" ]; then
    echo "No matching release found for track: $releaseTrack" >&2
    exit 1
  fi

  echo $latestVersion
}

# Function to get the Golang image hash.
getGolangImageHash() {
  local version=$1
  local tag=${version#go}
  local url="https://registry.hub.docker.com/v2/repositories/library/golang/tags/$tag"
  local images=$(curl -s $url)
  local imageHash=$(echo $images | jq -r '.digest' | head -n 1)

  if [ -z "$imageHash" ]; then
    echo "No amd64 image found for tag: $tag" >&2
    exit 1
  fi

  echo $imageHash
}

# Function to update Go version in files.
updateGoVersion() {
  local version=$1
  local imageHash=$2

  for filePath in "${filesToPatchWithGolangVersion[@]}"; do
    sed -i.bak -E "s/(BUILD_BASE_IMAGE[[:space:]]*\?=[[:space:]]*golang:)[0-9]+\.[0-9]+\.[0-9]+(@sha256:[a-f0-9]{64})?/\1${version#go}@${imageHash}/" $filePath
    sed -i.bak -E "s/(GO_VERSION:[[:space:]]*)[0-9]+\.[0-9]+\.[0-9]+/\1${version#go}/" $filePath
    rm "${filePath}.bak"
  done

  echo "Running go mod tidy to update generated files"
  go mod tidy
}

# Function to create a changelog entry.
create_changelog() {
  local targetVersion=$1
  local changelogFile="changelogs/unreleased/nnnn-$USER-small.md"
  local majorMinorVersion=$(echo $targetVersion | grep -oP 'go[0-9]+\.[0-9]+')
  local releaseNotesUrl="https://go.dev/doc/devel/release#${majorMinorVersion}.0"

  cat <<EOF > $changelogFile
Updates Go to $targetVersion. See the [Go release notes]($releaseNotesUrl) for more information about the content of the release.
EOF

  echo "Changelog created at $changelogFile"
}

# Main function.
main() {
  local releaseTrack=$1

  local latestVersion=$(getLatestGoVersion $releaseTrack)
  echo "Latest version: $latestVersion"

  local imageHash=$(getGolangImageHash $latestVersion)
  echo "Golang image hash: $imageHash"

  updateGoVersion $latestVersion $imageHash

  create_changelog $latestVersion
}

# Check if release track is provided.
if [[ -z ${1:-} ]]; then
  echo "Usage: $0 <releaseTrack>"
  exit 1
fi

main $1
