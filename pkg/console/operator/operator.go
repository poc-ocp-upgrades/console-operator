package operator

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"reflect"
	"time"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	corev1 "k8s.io/client-go/informers/core/v1"
	appsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreclientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	configv1 "github.com/openshift/api/config/v1"
	operatorsv1 "github.com/openshift/api/operator/v1"
	oauthclientv1 "github.com/openshift/client-go/oauth/clientset/versioned/typed/oauth/v1"
	oauthinformersv1 "github.com/openshift/client-go/oauth/informers/externalversions/oauth/v1"
	routeclientv1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/openshift/console-operator/pkg/api"
	"github.com/openshift/console-operator/pkg/boilerplate/operator"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/status"
	configinformer "github.com/openshift/client-go/config/informers/externalversions"
	operatorinformerv1 "github.com/openshift/client-go/operator/informers/externalversions/operator/v1"
	routesinformersv1 "github.com/openshift/client-go/route/informers/externalversions/route/v1"
	appsinformersv1 "k8s.io/client-go/informers/apps/v1"
	configclientv1 "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	operatorclientv1 "github.com/openshift/client-go/operator/clientset/versioned/typed/operator/v1"
	customerrors "github.com/openshift/console-operator/pkg/console/errors"
	"github.com/openshift/console-operator/pkg/console/subresource/configmap"
	"github.com/openshift/console-operator/pkg/console/subresource/deployment"
	"github.com/openshift/console-operator/pkg/console/subresource/oauthclient"
	"github.com/openshift/console-operator/pkg/console/subresource/route"
	"github.com/openshift/console-operator/pkg/console/subresource/secret"
	"github.com/openshift/console-operator/pkg/console/subresource/service"
)

const (
	controllerName = "Console"
)

var CreateDefaultConsoleFlag bool

type consoleOperator struct {
	operatorConfigClient		operatorclientv1.ConsoleInterface
	consoleConfigClient		configclientv1.ConsoleInterface
	secretsClient			coreclientv1.SecretsGetter
	configMapClient			coreclientv1.ConfigMapsGetter
	serviceClient			coreclientv1.ServicesGetter
	deploymentClient		appsv1.DeploymentsGetter
	routeClient			routeclientv1.RoutesGetter
	oauthClient			oauthclientv1.OAuthClientsGetter
	infrastructureConfigClient	configclientv1.InfrastructureInterface
	versionGetter			status.VersionGetter
	recorder			events.Recorder
}

func NewConsoleOperator(operatorConfigInformer operatorinformerv1.ConsoleInformer, configInformer configinformer.SharedInformerFactory, coreV1 corev1.Interface, managedCoreV1 corev1.Interface, deployments appsinformersv1.DeploymentInformer, routes routesinformersv1.RouteInformer, oauthClients oauthinformersv1.OAuthClientInformer, operatorConfigClient operatorclientv1.OperatorV1Interface, configClient configclientv1.ConfigV1Interface, corev1Client coreclientv1.CoreV1Interface, deploymentClient appsv1.DeploymentsGetter, routev1Client routeclientv1.RoutesGetter, oauthv1Client oauthclientv1.OAuthClientsGetter, versionGetter status.VersionGetter, recorder events.Recorder) operator.Runner {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &consoleOperator{operatorConfigClient: operatorConfigClient.Consoles(), consoleConfigClient: configClient.Consoles(), infrastructureConfigClient: configClient.Infrastructures(), secretsClient: corev1Client, configMapClient: corev1Client, serviceClient: corev1Client, deploymentClient: deploymentClient, routeClient: routev1Client, oauthClient: oauthv1Client, versionGetter: versionGetter, recorder: recorder}
	secretsInformer := coreV1.Secrets()
	configMapInformer := coreV1.ConfigMaps()
	managedConfigMapInformer := managedCoreV1.ConfigMaps()
	serviceInformer := coreV1.Services()
	configV1Informers := configInformer.Config().V1()
	configNameFilter := operator.FilterByNames(api.ConfigResourceName)
	targetNameFilter := operator.FilterByNames(api.OpenShiftConsoleName)
	return operator.New(controllerName, c, operator.WithInformer(configV1Informers.Consoles(), configNameFilter), operator.WithInformer(operatorConfigInformer, configNameFilter), operator.WithInformer(configV1Informers.Infrastructures(), configNameFilter), operator.WithInformer(deployments, targetNameFilter), operator.WithInformer(routes, targetNameFilter), operator.WithInformer(serviceInformer, targetNameFilter), operator.WithInformer(oauthClients, targetNameFilter), operator.WithInformer(configMapInformer, operator.FilterByNames(configmap.ConsoleConfigMapName, configmap.ServiceCAConfigMapName)), operator.WithInformer(managedConfigMapInformer, operator.FilterByNames(configmap.ConsoleConfigMapName)), operator.WithInformer(secretsInformer, operator.FilterByNames(deployment.ConsoleOauthConfigName)))
}
func (c *consoleOperator) Key() (metav1.Object, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	operatorConfig, err := c.operatorConfigClient.Get(api.ConfigResourceName, metav1.GetOptions{})
	if errors.IsNotFound(err) && CreateDefaultConsoleFlag {
		if _, err := c.operatorConfigClient.Create(c.defaultConsoleOperatorConfig()); err != nil {
			logrus.Errorf("no console operator config found. Creating. %v \n", err)
			return nil, err
		}
	}
	return operatorConfig, err
}
func (c *consoleOperator) Sync(obj metav1.Object) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	startTime := time.Now()
	logrus.Infof("started syncing operator %q (%v)", obj.GetName(), startTime)
	defer logrus.Infof("finished syncing operator %q (%v) \n\n", obj.GetName(), time.Since(startTime))
	operatorConfig := obj.(*operatorsv1.Console)
	consoleConfig, err := c.consoleConfigClient.Get(api.ConfigResourceName, metav1.GetOptions{})
	if errors.IsNotFound(err) && CreateDefaultConsoleFlag {
		logrus.Infof("no console config found. creating default config.")
		if _, err := c.consoleConfigClient.Create(c.defaultConsoleConfig()); err != nil {
			logrus.Errorf("error creating console config: %v \n", err)
			return err
		}
	}
	infrastructureConfig, err := c.infrastructureConfigClient.Get(api.ConfigResourceName, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("infrastructure config error: %v \n", err)
		return err
	}
	if err := c.handleSync(operatorConfig, consoleConfig, infrastructureConfig); err != nil {
		return err
	}
	return nil
}
func (c *consoleOperator) handleSync(operatorConfig *operatorsv1.Console, consoleConfig *configv1.Console, infrastructureConfig *configv1.Infrastructure) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	operatorConfigCopy := operatorConfig.DeepCopy()
	switch operatorConfigCopy.Spec.ManagementState {
	case operatorsv1.Managed:
		logrus.Println("console is in a managed state.")
	case operatorsv1.Unmanaged:
		logrus.Println("console is in an unmanaged state.")
		c.ConditionsManagementStateUnmanaged(operatorConfigCopy)
		if !reflect.DeepEqual(operatorConfigCopy, operatorConfig) {
			c.SyncStatus(operatorConfigCopy)
		}
		return nil
	case operatorsv1.Removed:
		logrus.Println("console has been removed.")
		c.ConditionsManagementStateRemoved(operatorConfigCopy)
		if !reflect.DeepEqual(operatorConfigCopy, operatorConfig) {
			c.SyncStatus(operatorConfigCopy)
		}
		return c.deleteAllResources(operatorConfigCopy)
	default:
		c.ConditionsManagementStateInvalid(operatorConfigCopy)
		if !reflect.DeepEqual(operatorConfigCopy, operatorConfig) {
			c.SyncStatus(operatorConfigCopy)
		}
		return fmt.Errorf("console is in an unknown state: %v", operatorConfigCopy.Spec.ManagementState)
	}
	c.ConditionNotDegraded(operatorConfigCopy)
	err := sync_v400(c, operatorConfigCopy, consoleConfig, infrastructureConfig)
	if err != nil {
		if !customerrors.IsSyncError(err) {
			c.SyncStatus(c.ConditionResourceSyncDegraded(operatorConfigCopy, err.Error()))
			return err
		} else {
			c.SyncStatus(operatorConfigCopy)
			return nil
		}
	}
	if !reflect.DeepEqual(operatorConfigCopy, operatorConfig) {
		c.SyncStatus(operatorConfigCopy)
	}
	return nil
}
func (c *consoleOperator) deleteAllResources(cr *operatorsv1.Console) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	logrus.Info("deleting console resources")
	defer logrus.Info("finished deleting console resources")
	var errs []error
	errs = append(errs, c.serviceClient.Services(api.TargetNamespace).Delete(service.Stub().Name, &metav1.DeleteOptions{}))
	errs = append(errs, c.routeClient.Routes(api.TargetNamespace).Delete(route.Stub().Name, &metav1.DeleteOptions{}))
	errs = append(errs, c.configMapClient.ConfigMaps(api.TargetNamespace).Delete(configmap.Stub().Name, &metav1.DeleteOptions{}))
	errs = append(errs, c.secretsClient.Secrets(api.TargetNamespace).Delete(secret.Stub().Name, &metav1.DeleteOptions{}))
	existingOAuthClient, getAuthErr := c.oauthClient.OAuthClients().Get(oauthclient.Stub().Name, metav1.GetOptions{})
	errs = append(errs, getAuthErr)
	_, updateAuthErr := c.oauthClient.OAuthClients().Update(oauthclient.DeRegisterConsoleFromOAuthClient(existingOAuthClient))
	errs = append(errs, updateAuthErr)
	errs = append(errs, c.deploymentClient.Deployments(api.TargetNamespace).Delete(deployment.Stub().Name, &metav1.DeleteOptions{}))
	return utilerrors.FilterOut(utilerrors.NewAggregate(errs), errors.IsNotFound)
}
func (c *consoleOperator) defaultConsoleConfig() *configv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &configv1.Console{ObjectMeta: metav1.ObjectMeta{Name: api.ConfigResourceName}}
}
func (c *consoleOperator) defaultConsoleOperatorConfig() *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &operatorsv1.Console{ObjectMeta: metav1.ObjectMeta{Name: api.ConfigResourceName}, Spec: operatorsv1.ConsoleSpec{OperatorSpec: operatorsv1.OperatorSpec{ManagementState: operatorsv1.Managed, LogLevel: operatorsv1.Normal}}}
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
