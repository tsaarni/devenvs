

# use major or major.minor as argument to pick the latest release from the release track
apps/update-envoy.sh 1.30
apps/update-go.sh 1

https://projectcontour.io/resources/release-process/

git remote add upstream git@github.com:projectcontour/contour.git

export CONTOUR_UPSTREAM_REMOTE_NAME=upstream



export CONTOUR_RELEASE_VERSION=v1.30.1
export CONTOUR_RELEASE_VERSION_MAJOR=1
export CONTOUR_RELEASE_VERSION_MINOR=30
export CONTOUR_PREVIOUS_VERSION=v1.30.0

git checkout release-1.30
git pull
git reset --hard
echo ./hack/release/make-release-tag.sh $CONTOUR_PREVIOUS_VERSION $CONTOUR_RELEASE_VERSION


git push ${CONTOUR_UPSTREAM_REMOTE_NAME} release-${CONTOUR_RELEASE_VERSION_MAJOR}.${CONTOUR_RELEASE_VERSION_MINOR}
git push ${CONTOUR_UPSTREAM_REMOTE_NAME} ${CONTOUR_RELEASE_VERSION}


docker pull ghcr.io/projectcontour/contour:${CONTOUR_RELEASE_VERSION}
docker run --rm  ghcr.io/projectcontour/contour:${CONTOUR_RELEASE_VERSION} contour serve



export CONTOUR_RELEASE_VERSION=v1.29.3
export CONTOUR_RELEASE_VERSION_MAJOR=1
export CONTOUR_RELEASE_VERSION_MINOR=29
export CONTOUR_PREVIOUS_VERSION=v1.29.2

git checkout release-1.29
git pull
git reset --hard
echo ./hack/release/make-release-tag.sh $CONTOUR_PREVIOUS_VERSION $CONTOUR_RELEASE_VERSION



git push ${CONTOUR_UPSTREAM_REMOTE_NAME} release-${CONTOUR_RELEASE_VERSION_MAJOR}.${CONTOUR_RELEASE_VERSION_MINOR}
git push ${CONTOUR_UPSTREAM_REMOTE_NAME} ${CONTOUR_RELEASE_VERSION}





export CONTOUR_RELEASE_VERSION=v1.28.7
export CONTOUR_RELEASE_VERSION_MAJOR=1
export CONTOUR_RELEASE_VERSION_MINOR=28
export CONTOUR_PREVIOUS_VERSION=v1.28.6

git checkout release-1.28
git pull
git reset --hard
echo ./hack/release/make-release-tag.sh $CONTOUR_PREVIOUS_VERSION $CONTOUR_RELEASE_VERSION


git push ${CONTOUR_UPSTREAM_REMOTE_NAME} release-${CONTOUR_RELEASE_VERSION_MAJOR}.${CONTOUR_RELEASE_VERSION_MINOR}
git push ${CONTOUR_UPSTREAM_REMOTE_NAME} ${CONTOUR_RELEASE_VERSION}





# envoy releases
https://github.com/envoyproxy/envoy/releases

# envoy images
https://hub.docker.com/r/envoyproxy/envoy/tags?name=v1.29.10

# go releases
https://go.dev/doc/devel/release


# patch releases


# release-1.30: Bump Envoy to v1.31.2 and Go to v1.22.8
https://github.com/projectcontour/contour/pull/6715
# release-1.29: Bump Envoy to v1.30.6 and Go to v1.22.8
https://github.com/projectcontour/contour/pull/6716
# release-1.28: Bump Envoy to v1.29.9 and Go to v1.21.13
https://github.com/projectcontour/contour/pull/6717

# Patch release updates
https://github.com/projectcontour/contour/pull/6724


git checkout release-1.30-envoy-bump
git checkout release-1.29-maint-bumps
git checkout release-1.28-maint-bumps


c Makefile
c cmd/contour/gatewayprovisioner.go
c 03-envoy.yaml
c examples/contour/03-envoy.yaml
c examples/deployment/03-envoy-deployment.yaml
make generate


git status

        modified:   Makefile
        modified:   cmd/contour/gatewayprovisioner.go
        modified:   examples/contour/03-envoy.yaml
        modified:   examples/deployment/03-envoy-deployment.yaml
        modified:   examples/render/contour-deployment.yaml
        modified:   examples/render/contour-gateway.yaml
        modified:   examples/render/contour.yaml

git add -u
git commit -sm "Bump Envoy to v1.29.10"
git push nordix





git checkout release-updates


c changelogs/CHANGELOG-v1.28.7.md
c changelogs/CHANGELOG-v1.29.3.md
c changelogs/CHANGELOG-v1.30.1.md
c site/content/resources/compatibility-matrix.md
c versions.yaml



# Test after release

kind delete cluster --name contour
kind create cluster --config ~/work/devenvs/contour/configs/kind-cluster-config.yaml --name contour
kubectl apply -f https://projectcontour.io/quickstart/contour.yaml

# check that pod is up
kubectl -n projectcontour get pod

# check image version
kubectl -n projectcontour exec -it $(kubectl -n projectcontour get pod -l app=contour -o jsonpath='{.items[0].metadata.name}') -- contour version

# check that traffic works
kubectl apply -f manifests/echoserver.yaml
http http://echoserver.127-0-0-101.nip.io



# Update slack channel title by clicking "edit"

Contour 1.30.1 out now! :tada: https://github.com/projectcontour/contour/releases/tag/v1.30.1 ||
https://github.com/projectcontour/contour/blob/main/changelogs/CHANGELOG-v1.30.1.md || https://github.com/projectcontour/contour



# Post on the channel
New Contour patch versions 1.30.1, 1.29.3 and 1.28.7 have been released! They include the latest Envoy releases and other dependency updates.


https://groups.google.com/g/project-contour

# Pick project-contour as the sender


Subject:
Contour v1.30.1, v1.29.3, and v1.28.7 released

Hi,

New Contour patch versions v1.30.1, v1.29.3 and v1.28.7 have been released! They include the latest Envoy releases and other dependency updates.

* v1.30.1 https://github.com/projectcontour/contour/releases/tag/v1.30.1
* v1.29.3 https://github.com/projectcontour/contour/releases/tag/v1.29.3
* v1.28.7 https://github.com/projectcontour/contour/releases/tag/v1.28.7


Thanks,
The Contour team
