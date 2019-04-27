package operator

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"github.com/sirupsen/logrus"
	oauthv1 "github.com/openshift/api/oauth/v1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	configv1 "github.com/openshift/api/config/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/console-operator/pkg/api"
	"github.com/openshift/console-operator/pkg/console/subresource/util"
	"github.com/openshift/console-operator/pkg/crypto"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resource/resourceapply"
	"github.com/openshift/library-go/pkg/operator/resource/resourcemerge"
	customerrors "github.com/openshift/console-operator/pkg/console/errors"
	configmapsub "github.com/openshift/console-operator/pkg/console/subresource/configmap"
	deploymentsub "github.com/openshift/console-operator/pkg/console/subresource/deployment"
	oauthsub "github.com/openshift/console-operator/pkg/console/subresource/oauthclient"
	routesub "github.com/openshift/console-operator/pkg/console/subresource/route"
	secretsub "github.com/openshift/console-operator/pkg/console/subresource/secret"
	servicesub "github.com/openshift/console-operator/pkg/console/subresource/service"
)

func sync_v400(co *consoleOperator, operatorConfig *operatorv1.Console, consoleConfig *configv1.Console, infrastructureConfig *configv1.Infrastructure) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	logrus.Println("running sync loop 4.0.0")
	recorder := co.recorder
	toUpdate := false
	rt, rtChanged, rtErr := SyncRoute(co, operatorConfig)
	if rtErr != nil {
		msg := fmt.Sprintf("%v: %s\n", "route", rtErr)
		logrus.Printf("incomplete sync: %v \n", msg)
		co.ConditionResourceSyncProgressing(operatorConfig, msg)
		return rtErr
	}
	toUpdate = toUpdate || rtChanged
	svc, svcChanged, svcErr := SyncService(co, recorder, operatorConfig)
	if svcErr != nil {
		msg := fmt.Sprintf("%q: %v\n", "service", svcErr)
		logrus.Printf("incomplete sync: %v \n", msg)
		co.ConditionResourceSyncProgressing(operatorConfig, msg)
		return svcErr
	}
	toUpdate = toUpdate || svcChanged
	cm, cmChanged, cmErr := SyncConfigMap(co, recorder, operatorConfig, consoleConfig, infrastructureConfig, rt)
	if cmErr != nil {
		msg := fmt.Sprintf("%q: %v\n", "configmap", cmErr)
		logrus.Printf("incomplete sync: %v \n", msg)
		co.ConditionResourceSyncProgressing(operatorConfig, msg)
		return cmErr
	}
	toUpdate = toUpdate || cmChanged
	serviceCAConfigMap, serviceCAConfigMapChanged, serviceCAConfigMapErr := SyncServiceCAConfigMap(co, operatorConfig)
	if serviceCAConfigMapErr != nil {
		msg := fmt.Sprintf("%q: %v\n", "serviceCAconfigmap", serviceCAConfigMapErr)
		logrus.Printf("incomplete sync: %v \n", msg)
		co.ConditionResourceSyncProgressing(operatorConfig, msg)
		return serviceCAConfigMapErr
	}
	toUpdate = toUpdate || serviceCAConfigMapChanged
	sec, secChanged, secErr := SyncSecret(co, recorder, operatorConfig)
	if secErr != nil {
		msg := fmt.Sprintf("%q: %v\n", "secret", secErr)
		logrus.Printf("incomplete sync: %v \n", msg)
		co.ConditionResourceSyncProgressing(operatorConfig, msg)
		return secErr
	}
	toUpdate = toUpdate || secChanged
	oauthClient, oauthChanged, oauthErr := SyncOAuthClient(co, operatorConfig, sec, rt)
	if oauthErr != nil {
		msg := fmt.Sprintf("%q: %v\n", "oauth", oauthErr)
		logrus.Printf("incomplete sync: %v \n", msg)
		co.ConditionResourceSyncProgressing(operatorConfig, msg)
		return oauthErr
	}
	toUpdate = toUpdate || oauthChanged
	actualDeployment, depChanged, depErr := SyncDeployment(co, recorder, operatorConfig, cm, serviceCAConfigMap, sec, rt)
	if depErr != nil {
		msg := fmt.Sprintf("%q: %v\n", "deployment", depErr)
		logrus.Printf("incomplete sync: %v \n", msg)
		co.ConditionResourceSyncProgressing(operatorConfig, msg)
		return depErr
	}
	toUpdate = toUpdate || depChanged
	resourcemerge.SetDeploymentGeneration(&operatorConfig.Status.Generations, actualDeployment)
	operatorConfig.Status.ObservedGeneration = operatorConfig.ObjectMeta.Generation
	logrus.Println("-----------------------")
	logrus.Printf("sync loop 4.0.0 resources updated: %v \n", toUpdate)
	logrus.Println("-----------------------")
	co.ConditionResourceSyncSuccess(operatorConfig)
	if toUpdate {
		co.ConditionResourceSyncProgressing(operatorConfig, "Changes made during sync updates, additional sync expected.")
	} else {
		version := os.Getenv("RELEASE_VERSION")
		if !deploymentsub.IsAvailableAndUpdated(actualDeployment) {
			co.ConditionResourceSyncProgressing(operatorConfig, fmt.Sprintf("Moving to version %s", strings.Split(version, "-")[0]))
		} else {
			if co.versionGetter.GetVersions()["operator"] != version {
				co.versionGetter.SetVersion("operator", version)
			}
			co.ConditionResourceSyncNotProgressing(operatorConfig)
		}
	}
	if !deploymentsub.IsReady(actualDeployment) {
		msg := fmt.Sprintf("%v pods available for console deployment", actualDeployment.Status.ReadyReplicas)
		logrus.Println(msg)
		co.ConditionDeploymentNotAvailable(operatorConfig, msg)
	} else if !routesub.IsAdmitted(rt) {
		logrus.Println("console route is not admitted")
		co.SetStatusCondition(operatorConfig, operatorv1.OperatorStatusTypeAvailable, operatorv1.ConditionFalse, "RouteNotAdmitted", "console route is not admitted")
	} else {
		co.ConditionDeploymentAvailable(operatorConfig)
	}
	logrus.Println("sync_v400: updating console status")
	if _, err := SyncConsoleConfig(co, consoleConfig, rt); err != nil {
		logrus.Errorf("could not update console config status: %v \n", err)
		return err
	}
	defer func() {
		logrus.Printf("sync loop 4.0.0 complete")
		if svcChanged {
			logrus.Printf("\t service changed: %v", svc.GetResourceVersion())
		}
		if rtChanged {
			logrus.Printf("\t route changed: %v", rt.GetResourceVersion())
		}
		if cmChanged {
			logrus.Printf("\t configmap changed: %v", cm.GetResourceVersion())
		}
		if serviceCAConfigMapChanged {
			logrus.Printf("\t service-ca configmap changed: %v", serviceCAConfigMap.GetResourceVersion())
		}
		if secChanged {
			logrus.Printf("\t secret changed: %v", sec.GetResourceVersion())
		}
		if oauthChanged {
			logrus.Printf("\t oauth changed: %v", oauthClient.GetResourceVersion())
		}
		if depChanged {
			logrus.Printf("\t deployment changed: %v", actualDeployment.GetResourceVersion())
		}
	}()
	return nil
}
func SyncConsoleConfig(co *consoleOperator, consoleConfig *configv1.Console, route *routev1.Route) (*configv1.Console, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	host := routesub.GetCanonicalHost(route)
	if host == "" {
		customErr := customerrors.NewSyncError("waiting on route host")
		logrus.Errorf("%q: %v \n", "route", customErr)
		return nil, customErr
	}
	httpsHost := util.HTTPS(host)
	if consoleConfig.Status.ConsoleURL != httpsHost {
		logrus.Printf("updating console.config.openshift.io with hostname: %v \n", host)
		consoleConfig.Status.ConsoleURL = httpsHost
	}
	return co.consoleConfigClient.UpdateStatus(consoleConfig)
}
func SyncDeployment(co *consoleOperator, recorder events.Recorder, operatorConfig *operatorv1.Console, cm *corev1.ConfigMap, serviceCAConfigMap *corev1.ConfigMap, sec *corev1.Secret, rt *routev1.Route) (*appsv1.Deployment, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	requiredDeployment := deploymentsub.DefaultDeployment(operatorConfig, cm, serviceCAConfigMap, sec, rt)
	expectedGeneration := getDeploymentGeneration(co)
	genChanged := operatorConfig.ObjectMeta.Generation != operatorConfig.Status.ObservedGeneration
	if genChanged {
		logrus.Printf("deployment generation changed from %s to %s \n", operatorConfig.ObjectMeta.Generation, operatorConfig.Status.ObservedGeneration)
	}
	deploymentsub.LogDeploymentAnnotationChanges(co.deploymentClient, requiredDeployment)
	deployment, deploymentChanged, applyDepErr := resourceapply.ApplyDeployment(co.deploymentClient, recorder, requiredDeployment, expectedGeneration, genChanged)
	if applyDepErr != nil {
		logrus.Errorf("%q: %v \n", "deployment", applyDepErr)
		return nil, false, applyDepErr
	}
	return deployment, deploymentChanged, nil
}
func SyncOAuthClient(co *consoleOperator, operatorConfig *operatorv1.Console, sec *corev1.Secret, rt *routev1.Route) (*oauthv1.OAuthClient, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	host := routesub.GetCanonicalHost(rt)
	if host == "" {
		customErr := customerrors.NewSyncError("waiting on route host")
		logrus.Errorf("%q: %v \n", "oauth", customErr)
		return nil, false, customErr
	}
	oauthClient, err := co.oauthClient.OAuthClients().Get(oauthsub.Stub().Name, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("%q: %v \n", "oauth", err)
		return nil, false, errors.New("oauth client for console does not exist and cannot be created")
	}
	oauthsub.RegisterConsoleToOAuthClient(oauthClient, host, secretsub.GetSecretString(sec))
	oauthClient, oauthChanged, oauthErr := oauthsub.CustomApplyOAuth(co.oauthClient, oauthClient)
	if oauthErr != nil {
		logrus.Errorf("%q: %v \n", "oauth", oauthErr)
		return nil, false, oauthErr
	}
	return oauthClient, oauthChanged, nil
}
func SyncSecret(co *consoleOperator, recorder events.Recorder, operatorConfig *operatorv1.Console) (*corev1.Secret, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	secret, err := co.secretsClient.Secrets(api.TargetNamespace).Get(secretsub.Stub().Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) || secretsub.GetSecretString(secret) == "" {
		return resourceapply.ApplySecret(co.secretsClient, recorder, secretsub.DefaultSecret(operatorConfig, crypto.Random256BitsString()))
	}
	if err != nil {
		logrus.Errorf("%q: %v \n", "secret", err)
		return nil, false, err
	}
	return secret, false, nil
}
func SyncConfigMap(co *consoleOperator, recorder events.Recorder, operatorConfig *operatorv1.Console, consoleConfig *configv1.Console, infrastructureConfig *configv1.Infrastructure, rt *routev1.Route) (*corev1.ConfigMap, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	managedConfig, mcErr := co.configMapClient.ConfigMaps(api.OpenshiftConfigManagedNamespace).Get(api.OpenShiftConsoleConfigMapName, metav1.GetOptions{})
	if mcErr != nil && !apierrors.IsNotFound(mcErr) {
		logrus.Errorf("managed config error: %v \n", mcErr)
		return nil, false, mcErr
	}
	defaultConfigmap, _, err := configmapsub.DefaultConfigMap(operatorConfig, consoleConfig, managedConfig, infrastructureConfig, rt)
	if err != nil {
		return nil, false, err
	}
	cm, cmChanged, cmErr := resourceapply.ApplyConfigMap(co.configMapClient, recorder, defaultConfigmap)
	if cmErr != nil {
		logrus.Errorf("%q: %v \n", "configmap", cmErr)
		return nil, false, cmErr
	}
	if cmChanged {
		logrus.Println("new console config yaml:")
		logrus.Printf("%s \n", cm.Data)
	}
	return cm, cmChanged, cmErr
}
func SyncServiceCAConfigMap(co *consoleOperator, operatorConfig *operatorv1.Console) (*corev1.ConfigMap, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	required := configmapsub.DefaultServiceCAConfigMap(operatorConfig)
	existing, err := co.configMapClient.ConfigMaps(required.Namespace).Get(required.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := co.configMapClient.ConfigMaps(required.Namespace).Create(required)
		if err == nil {
			logrus.Println("service-ca configmap created")
		} else {
			logrus.Errorf("%q: %v \n", "service-ca configmap", err)
		}
		return actual, true, err
	}
	if err != nil {
		logrus.Errorf("%q: %v \n", "service-ca configmap", err)
		return nil, false, err
	}
	modified := resourcemerge.BoolPtr(false)
	resourcemerge.EnsureObjectMeta(modified, &existing.ObjectMeta, required.ObjectMeta)
	if !*modified {
		logrus.Println("service-ca configmap exists and is in the correct state")
		return existing, false, nil
	}
	actual, err := co.configMapClient.ConfigMaps(required.Namespace).Update(existing)
	if err == nil {
		logrus.Println("service-ca configmap updated")
	} else {
		logrus.Errorf("%q: %v \n", "service-ca configmap", err)
	}
	return actual, true, err
}
func SyncService(co *consoleOperator, recorder events.Recorder, operatorConfig *operatorv1.Console) (*corev1.Service, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	svc, svcChanged, svcErr := resourceapply.ApplyService(co.serviceClient, recorder, servicesub.DefaultService(operatorConfig))
	if svcErr != nil {
		logrus.Errorf("%q: %v \n", "service", svcErr)
		return nil, false, svcErr
	}
	return svc, svcChanged, svcErr
}
func SyncRoute(co *consoleOperator, operatorConfig *operatorv1.Console) (*routev1.Route, bool, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	rt, rtIsNew, rtErr := routesub.GetOrCreate(co.routeClient, routesub.DefaultRoute(operatorConfig))
	if rtErr != nil {
		logrus.Errorf("%q: %v \n", "route", rtErr)
		return nil, false, rtErr
	}
	host := routesub.GetCanonicalHost(rt)
	if host == "" {
		customErr := customerrors.NewSyncError("waiting on route host")
		logrus.Errorf("%q: %v \n", "route", customErr)
		return nil, false, customErr
	}
	if validatedRoute, changed := routesub.Validate(rt); changed {
		if _, err := co.routeClient.Routes(api.TargetNamespace).Update(validatedRoute); err != nil {
			logrus.Errorf("%q: %v \n", "route", err)
			return nil, false, err
		}
		customErr := customerrors.NewSyncError("route is invalid, correcting route state")
		logrus.Error(customErr)
		return nil, true, customErr
	}
	return rt, rtIsNew, rtErr
}
func getDeploymentGeneration(co *consoleOperator) int64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	deployment, err := co.deploymentClient.Deployments(api.TargetNamespace).Get(deploymentsub.Stub().Name, metav1.GetOptions{})
	if err != nil {
		return -1
	}
	return deployment.Generation
}
