package configmap

import (
	corev1 "k8s.io/api/core/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/console-operator/pkg/console/subresource/util"
)

const (
	ServiceCAConfigMapName		= "service-ca"
	injectCABundleAnnotation	= "service.alpha.openshift.io/inject-cabundle"
)

func DefaultServiceCAConfigMap(cr *operatorv1.Console) *corev1.ConfigMap {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	configMap := ServiceCAStub()
	util.AddOwnerRef(configMap, util.OwnerRefFrom(cr))
	return configMap
}
func ServiceCAStub() *corev1.ConfigMap {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	meta := util.SharedMeta()
	meta.Name = ServiceCAConfigMapName
	meta.Annotations = map[string]string{injectCABundleAnnotation: "true"}
	configMap := &corev1.ConfigMap{ObjectMeta: meta}
	return configMap
}
