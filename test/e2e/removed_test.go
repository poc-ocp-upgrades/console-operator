package e2e

import (
	"testing"
	"github.com/openshift/console-operator/pkg/testframework"
)

func TestRemoved(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client := testframework.MustNewClientset(t, nil)
	defer testframework.MustManageConsole(t, client)
	testframework.MustRemoveConsole(t, client)
	t.Logf("validating that the operator does not recreate removed resources when ManagementState:Removed...")
	errChan := make(chan error)
	go testframework.IsResourceUnavailable(errChan, client, "ConfigMap")
	go testframework.IsResourceUnavailable(errChan, client, "Route")
	go testframework.IsResourceUnavailable(errChan, client, "Service")
	go testframework.IsResourceUnavailable(errChan, client, "Deployment")
	checkErr := <-errChan
	if checkErr != nil {
		t.Fatal(checkErr)
	}
}
