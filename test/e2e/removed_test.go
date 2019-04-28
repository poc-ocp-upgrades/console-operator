package e2e

import (
	"testing"
	"github.com/openshift/console-operator/pkg/testframework"
)

func TestRemoved(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client := testframework.MustNewClientset(t, nil)
	defer testframework.MustManageConsole(t, client)
	testframework.MustRemoveConsole(t, client)
	t.Logf("validating that the operator does not recreate removed resources when ManagementState:Removed...")
	err := testframework.ConsoleResourcesUnavailable(client)
	if err != nil {
		t.Fatal(err)
	}
}
