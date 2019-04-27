package testframework

import (
	"fmt"
	"testing"
	"time"
	routev1 "github.com/openshift/api/route/v1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	consoleapi "github.com/openshift/console-operator/pkg/api"
)

var (
	AsyncOperationTimeout = 5 * time.Minute
)

func DeleteAll(t *testing.T, client *Clientset) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	resources := []string{"Deployment", "Service", "Route", "ConfigMap"}
	for _, resource := range resources {
		t.Logf("deleting console %s...", resource)
		if err := DeleteCompletely(func() (runtime.Object, error) {
			return GetResource(client, resource)
		}, func(*metav1.DeleteOptions) error {
			return deleteResource(client, resource)
		}); err != nil {
			t.Fatalf("unable to delete console %s: %s", resource, err)
		}
	}
}
func GetResource(client *Clientset, resource string) (runtime.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var res runtime.Object
	var err error
	switch resource {
	case "ConfigMap":
		res, err = GetConsoleConfigMap(client)
	case "Service":
		res, err = GetConsoleService(client)
	case "Route":
		res, err = GetConsoleRoute(client)
	case "Deployment":
		fallthrough
	default:
		res, err = GetConsoleDeployment(client)
	}
	return res, err
}
func GetConsoleConfigMap(client *Clientset) (*corev1.ConfigMap, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return client.ConfigMaps(consoleapi.OpenShiftConsoleNamespace).Get(consoleapi.OpenShiftConsoleConfigMapName, metav1.GetOptions{})
}
func GetConsoleService(client *Clientset) (*corev1.Service, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return client.Services(consoleapi.OpenShiftConsoleNamespace).Get(consoleapi.OpenShiftConsoleServiceName, metav1.GetOptions{})
}
func GetConsoleRoute(client *Clientset) (*routev1.Route, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return client.Routes(consoleapi.OpenShiftConsoleNamespace).Get(consoleapi.OpenShiftConsoleRouteName, metav1.GetOptions{})
}
func GetConsoleDeployment(client *Clientset) (*appv1.Deployment, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return client.Deployments(consoleapi.OpenShiftConsoleNamespace).Get(consoleapi.OpenShiftConsoleDeploymentName, metav1.GetOptions{})
}
func deleteResource(client *Clientset, resource string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var err error
	switch resource {
	case "ConfigMap":
		err = client.ConfigMaps(consoleapi.OpenShiftConsoleNamespace).Delete(consoleapi.OpenShiftConsoleConfigMapName, &metav1.DeleteOptions{})
	case "Service":
		err = client.Services(consoleapi.OpenShiftConsoleNamespace).Delete(consoleapi.OpenShiftConsoleServiceName, &metav1.DeleteOptions{})
	case "Route":
		err = client.Routes(consoleapi.OpenShiftConsoleNamespace).Delete(consoleapi.OpenShiftConsoleRouteName, &metav1.DeleteOptions{})
	case "Deployment":
		fallthrough
	default:
		err = client.Deployments(consoleapi.OpenShiftConsoleNamespace).Delete(consoleapi.OpenShiftConsoleDeploymentName, &metav1.DeleteOptions{})
	}
	return err
}
func DeleteCompletely(getObject func() (runtime.Object, error), deleteObject func(*metav1.DeleteOptions) error) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	obj, err := getObject()
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	accessor, err := meta.Accessor(obj)
	uid := accessor.GetUID()
	policy := metav1.DeletePropagationForeground
	if err := deleteObject(&metav1.DeleteOptions{Preconditions: &metav1.Preconditions{UID: &uid}, PropagationPolicy: &policy}); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return wait.Poll(1*time.Second, AsyncOperationTimeout, func() (stop bool, err error) {
		obj, err = getObject()
		if err != nil {
			if errors.IsNotFound(err) {
				return true, nil
			}
			return false, err
		}
		accessor, err := meta.Accessor(obj)
		return accessor.GetUID() != uid, nil
	})
}
func IsResourceAvailable(errChan chan error, client *Clientset, resource string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	counter := 0
	maxCount := 30
	err := wait.Poll(1*time.Second, AsyncOperationTimeout, func() (stop bool, err error) {
		_, err = GetResource(client, resource)
		if err == nil {
			return true, nil
		}
		if counter == maxCount {
			if err != nil {
				return true, fmt.Errorf("deleted console %s was not recreated", resource)
			}
			return true, nil
		}
		counter++
		return false, nil
	})
	errChan <- err
}
func IsResourceUnavailable(errChan chan error, client *Clientset, resource string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	counter := 0
	maxCount := 15
	err := wait.Poll(1*time.Second, AsyncOperationTimeout, func() (stop bool, err error) {
		_, err = GetResource(client, resource)
		if err == nil {
			return true, fmt.Errorf("deleted console %s was recreated", resource)
		}
		if !errors.IsNotFound(err) {
			return true, err
		}
		counter++
		if counter == maxCount {
			return true, nil
		}
		return false, nil
	})
	errChan <- err
}
