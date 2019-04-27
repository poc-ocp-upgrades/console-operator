package testframework

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"testing"
	"github.com/davecgh/go-spew/spew"
	consoleapi "github.com/openshift/console-operator/pkg/api"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodLog []string

func (log PodLog) Contains(re *regexp.Regexp) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, line := range log {
		if re.MatchString(line) {
			return true
		}
	}
	return false
}

type PodSetLogs map[string]PodLog

func (psl PodSetLogs) Contains(re *regexp.Regexp) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, podlog := range psl {
		if podlog.Contains(re) {
			return true
		}
	}
	return false
}
func GetLogsByLabelSelector(client *Clientset, namespace string, labelSelector *metav1.LabelSelector) (PodSetLogs, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return nil, err
	}
	podList, err := client.Pods(namespace).List(metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}
	podLogs := make(PodSetLogs)
	for _, pod := range podList.Items {
		var podLog PodLog
		log, err := client.Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{}).Stream()
		if err != nil {
			return nil, fmt.Errorf("failed to get logs for pod %s: %s", pod.Name, err)
		}
		r := bufio.NewReader(log)
		for {
			line, readErr := r.ReadSlice('\n')
			if len(line) > 0 || readErr == nil {
				podLog = append(podLog, string(line))
			}
			if readErr == io.EOF {
				break
			} else if readErr != nil {
				return nil, fmt.Errorf("failed to read log for pod %s: %s", pod.Name, readErr)
			}
		}
		podLogs[pod.Name] = podLog
	}
	return podLogs, nil
}
func DumpObject(t *testing.T, prefix string, obj interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	t.Logf("%s:\n%s", prefix, spew.Sdump(obj))
}
func DumpPodLogs(t *testing.T, podLogs PodSetLogs) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(podLogs) > 0 {
		for pod, logs := range podLogs {
			t.Logf("=== logs for pod/%s", pod)
			for _, line := range logs {
				t.Logf("%s", line)
			}
		}
		t.Logf("=== end of logs")
	}
}
func GetOperatorLogs(client *Clientset) (PodSetLogs, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return GetLogsByLabelSelector(client, consoleapi.OpenShiftConsoleNamespace, &metav1.LabelSelector{MatchLabels: map[string]string{"name": "console-operator"}})
}
func DumpOperatorLogs(t *testing.T, client *Clientset) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	podLogs, err := GetOperatorLogs(client)
	if err != nil {
		t.Logf("failed to get the operator logs: %s", err)
	}
	DumpPodLogs(t, podLogs)
}
