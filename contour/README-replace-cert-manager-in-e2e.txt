
cd ~/work/contour-worktree/e2e-test-remove-cert-manager

make e2e


or run just single test

make setup-kind-cluster load-contour-image-kind 
CONTOUR_E2E_TEST_FOCUS="external name services work over https" make run-e2e
make cleanup-kind


