#! /usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

# Define the files to patch with the Envoy version.
filesToPatchWithEnvoyVersion=(
  "Makefile"
  "cmd/contour/gatewayprovisioner.go"
  "examples/contour/03-envoy.yaml"
  "examples/deployment/03-envoy-deployment.yaml"
)

# Function to get the latest Envoy tag.
get_latest_envoy_tag() {
  local releaseTrack=$1
  local url="https://api.github.com/repos/envoyproxy/envoy/releases"
  local releases=$(curl -s $url | jq -r ".[].tag_name")

  for tag in $releases; do
    if [[ $tag == v$releaseTrack.* ]]; then
      echo $tag
      return 0
    fi
  done

  echo "No matching release found for track: $releaseTrack" >&2
  return 1
}

# Function to update the Envoy image in the specified files.
update_envoy_image() {
  local targetVersion=$1
  local pattern="docker.io/envoyproxy/envoy:v[0-9]+\.[0-9]+\.[0-9]+"

  for filePath in "${filesToPatchWithEnvoyVersion[@]}"; do
    if [[ -f $filePath ]]; then
      sed -i.bak -E "s|$pattern|docker.io/envoyproxy/envoy:$targetVersion|g" $filePath
      rm "${filePath}.bak"
    else
      echo "File not found: $filePath" >&2
      return 1
    fi
  done

  echo "Running make generate to update generated files"
  make generate
}

# Function to create a changelog entry.
create_changelog() {
  local targetVersion=$1
  local changelogFile="changelogs/unreleased/nnnn-$USER-small.md"
  local majorMinorVersion=$(echo $targetVersion | grep -oP 'v[0-9]+\.[0-9]+')
  local releaseNotesUrl="https://www.envoyproxy.io/docs/envoy/$targetVersion/version_history/$majorMinorVersion/$targetVersion"

  cat <<EOF > $changelogFile
Updates Envoy to $targetVersion. See the [Envoy release notes]($releaseNotesUrl) for more information about the content of the release.
EOF

  echo "Changelog created at $changelogFile"
}

# Main function.
main() {
  local releaseTrack=$1
  local targetVersion=$(get_latest_envoy_tag $releaseTrack)

  if [[ $? -ne 0 ]]; then
    echo "Error getting latest version" >&2
    return 1
  fi

  if [[ -z $targetVersion ]]; then
    echo "No valid target version found" >&2
    return 1
  fi

  echo "Latest version: $targetVersion"

  update_envoy_image $targetVersion

  create_changelog $targetVersion

  echo "Update files (in main branch):"
  echo "  site/content/resources/compatibility-matrix.md"
  echo "  versions.yaml"
  echo "  changelog (if doing bump in main branch)"
}

# Check if release track is provided
if [[ -z ${1:-} ]]; then
  echo "Usage: $0 <releaseTrack>"
  exit 1
fi

main $1
