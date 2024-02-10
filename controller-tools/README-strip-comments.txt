
# run tests (from test.sh)
go test -race ./pkg/... ./cmd/... -parallel 4


# re-generate testdata according to changed controller-gen implementation
cd pkg/crd/testdata/
go generate
