package e2e

import (
	"reflect"
	"testing"
	"time"
	"github.com/openshift/console-operator/pkg/testframework"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	consoleapi "github.com/openshift/console-operator/pkg/api"
)

var pollTimeout = 10 * time.Second

func patchAndCheckConfigMap(t *testing.T, client *testframework.Clientset, isOperatorManaged bool) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Logf("patching Data on the console ConfigMap")
	configMap, err := client.ConfigMaps(consoleapi.OpenShiftConsoleNamespace).Patch(consoleapi.OpenShiftConsoleConfigMapName, types.MergePatchType, []byte(`{"data": {"console-config.yaml": "test"}}`))
	if err != nil {
		return err
	}
	patchedData := configMap.Data
	t.Logf("polling for patched Data on the console ConfigMap")
	err = wait.Poll(1*time.Second, pollTimeout, func() (stop bool, err error) {
		configMap, err = testframework.GetConsoleConfigMap(client)
		if err != nil {
			return true, err
		}
		newData := configMap.Data
		if isOperatorManaged {
			return !reflect.DeepEqual(patchedData, newData), nil
		}
		return reflect.DeepEqual(patchedData, newData), nil
	})
	return err
}
func patchAndCheckService(t *testing.T, client *testframework.Clientset, isOperatorManaged bool) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Logf("patching Annotation on the console Service")
	service, err := client.Services(consoleapi.OpenShiftConsoleNamespace).Patch(consoleapi.OpenShiftConsoleServiceName, types.MergePatchType, []byte(`{"metadata": {"annotations": {"service.alpha.openshift.io/serving-cert-secret-name": "test"}}}`))
	if err != nil {
		return err
	}
	patchedData := service.GetAnnotations()
	t.Logf("polling for patched Annotation on the console Service")
	err = wait.Poll(1*time.Second, pollTimeout, func() (stop bool, err error) {
		service, err = testframework.GetConsoleService(client)
		if err != nil {
			return true, err
		}
		newData := service.GetAnnotations()
		if isOperatorManaged {
			return !reflect.DeepEqual(patchedData, newData), nil
		}
		return reflect.DeepEqual(patchedData, newData), nil
	})
	return err
}
func patchAndCheckRoute(t *testing.T, client *testframework.Clientset, isOperatorManaged bool) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Logf("patching TargetPort on the console Route")
	route, err := client.Routes(consoleapi.OpenShiftConsoleNamespace).Patch(consoleapi.OpenShiftConsoleRouteName, types.MergePatchType, []byte(`{"spec": {"port": {"targetPort": "http"}}}`))
	if err != nil {
		return err
	}
	patchedData := route.Spec.Port.TargetPort
	t.Logf("polling for patched TargetPort on the console Route")
	err = wait.Poll(1*time.Second, pollTimeout, func() (stop bool, err error) {
		route, err = testframework.GetConsoleRoute(client)
		if err != nil {
			return true, err
		}
		newData := route.Spec.Port.TargetPort
		if isOperatorManaged {
			return !reflect.DeepEqual(patchedData, newData), nil
		}
		return reflect.DeepEqual(patchedData, newData), nil
	})
	return err
}
