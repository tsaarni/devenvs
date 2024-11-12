# Go 2.23 crash with aliases

https://github.com/kubernetes-sigs/controller-tools/pull/1079
https://github.com/kubernetes-sigs/controller-tools/pull/1078
https://github.com/kubernetes-sigs/controller-tools/pull/1061



# Crash seems to be triggered by

ContourConfiguration.Status.Conditions[].Condition

type Condition = meta_v1.Condition



### Test procedure

# add new structs to test aliases

pkg/crd/testdata/cronjob_types.go

cd ~/work/controller-tools/pkg/crd/testdata
go generate


git diff   # check the changes


rm -rf foo
go run ./cmd/controller-gen  crd:crdVersions=v1 paths=/home/tsaarni/work/contour/apis/... output:dir=foo
foo/projectcontour.io_contourconfigurations.yaml


### Create binary



go build ./cmd/controller-gen


patch -p1 < ~/work/devenvs/controller-tools/run-local-controller-gen.patch


make generate-crd-yaml
bash -x ./hack/generate-crd-yaml.sh

make generate-crd-deepcopy
./hack/generate-crd-deepcopy.sh

make generate





/home/tsaarni/go/bin/gen-crd-api-reference-docs -template-dir /home/tsaarni/work/contour/hack/api-docs-config/refdocs -config /home/tsaarni/work/contour/hack/api-docs-config/refdocs/config.json -api-dir github.com/projectcontour/contour/apis/projectcontour -out-file /home/tsaarni/work/contour/site/content/docs/main/config/api-reference.html
I1109 08:45:46.473039 1320009

/home/tsaarni/work/gen-crd-api-reference-docs/gen-crd-api-reference-docs -template-dir /home/tsaarni/work/contour/hack/api-docs-config/refdocs -config /home/tsaarni/work/contour/hack/api-docs-config/refdocs/config.json -api-dir github.com/projectcontour/contour/apis/projectcontour -out-file /home/tsaarni/work/contour/site/content/docs/main/config/api-reference.html


go run github.com/ahmetb/gen-crd-api-reference-docs@v0.3.0 -template-dir /home/tsaarni/work/contour/hack/api-docs-config/refdocs -config /home/tsaarni/work/contour/hack/api-docs-config/refdocs/config.json -api-dir github.com/projectcontour/contour/apis/projectcontour -out-file /home/tsaarni/work/contour/site/content/docs/main/config/api-reference.html


# debugging info
go run -v -x


go run main.go -template-dir /home/tsaarni/work/contour/hack/api-docs-config/refdocs -config /home/tsaarni/work/contour/hack/api-docs-config/refdocs/config.json -api-dir github.com/projectcontour/contour/apis/projectcontour -out-file /home/tsaarni/work/contour/site/content/docs/main/config/api-reference.html



GODEBUG=gotypesalias=0 go run -x -v sigs.k8s.io/controller-tools/cmd/controller-gen crd:crdVersions=v1 paths=./apis/... output:dir=crd-tgW7rj





### Debug go generate in vscodee


mkdir -p .vscode
cp ~/work/devenvs/controller-tools/configs/launch.json .vscode


# breakpoints
# File: pkg/crd/schema.go
# Condition: (*ident).Name=="InlineAlias"
