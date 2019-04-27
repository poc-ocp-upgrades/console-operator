package operator

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	operatorsv1 "github.com/openshift/api/operator/v1"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
)

const (
	reasonAsExpected		= "AsExpected"
	reasonWorkloadFailing		= "WorkloadFailing"
	reasonUnmanaged			= "ManagementStateUnmanaged"
	reasonRemoved			= "ManagementStateRemoved"
	reasonInvalid			= "ManagementStateInvalid"
	reasonSyncLoopProgressing	= "SyncLoopProgressing"
	reasonSyncError			= "SynchronizationError"
	reasonNoPodsAvailable		= "NoPodsAvailable"
	reasonSyncLoopError		= "SyncLoopError"
)

func (c *consoleOperator) SyncStatus(operatorConfig *operatorsv1.Console) (*operatorsv1.Console, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.logConditions(operatorConfig.Status.Conditions)
	updatedConfig, err := c.operatorConfigClient.UpdateStatus(operatorConfig)
	if err != nil {
		errMsg := fmt.Errorf("status update error: %v \n", err)
		logrus.Error(errMsg)
		return nil, errMsg
	}
	return updatedConfig, nil
}
func (c *consoleOperator) logConditions(conditions []operatorsv1.OperatorCondition) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	logrus.Println("Operator.Status.Conditions")
	for _, condition := range conditions {
		buf := bytes.Buffer{}
		buf.WriteString(fmt.Sprintf("Status.Condition.%s: %s", condition.Type, condition.Status))
		hasMessage := condition.Message != ""
		hasReason := condition.Reason != ""
		if hasMessage && hasReason {
			buf.WriteString(" |")
			if hasReason {
				buf.WriteString(fmt.Sprintf(" (%s)", condition.Reason))
			}
			if hasMessage {
				buf.WriteString(fmt.Sprintf(" %s", condition.Message))
			}
		}
		logrus.Println(buf.String())
	}
}
func (c *consoleOperator) SetStatusCondition(operatorConfig *operatorsv1.Console, conditionType string, conditionStatus operatorsv1.ConditionStatus, conditionReason string, conditionMessage string) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: conditionType, Status: conditionStatus, Reason: conditionReason, Message: conditionMessage, LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionResourceSyncSuccess(operatorConfig *operatorsv1.Console) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeFailing, Status: operatorsv1.ConditionFalse, LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionDeploymentAvailable(operatorConfig *operatorsv1.Console) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeAvailable, Status: operatorsv1.ConditionTrue, LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionDeploymentNotAvailable(operatorConfig *operatorsv1.Console, message string) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeAvailable, Status: operatorsv1.ConditionFalse, Reason: reasonNoPodsAvailable, Message: message, LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionResourceSyncProgressing(operatorConfig *operatorsv1.Console, message string) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeProgressing, Status: operatorsv1.ConditionTrue, Reason: reasonSyncLoopProgressing, Message: message, LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionResourceSyncNotProgressing(operatorConfig *operatorsv1.Console) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeProgressing, Status: operatorsv1.ConditionFalse, LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionDegraded(operatorConfig *operatorsv1.Console, conditionReason string, conditionMessage string) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeDegraded, Status: operatorsv1.ConditionTrue, Reason: conditionReason, Message: conditionMessage, LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionNotDegraded(operatorConfig *operatorsv1.Console) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeDegraded, Status: operatorsv1.ConditionFalse, Reason: reasonAsExpected, LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionResourceSyncDegraded(operatorConfig *operatorsv1.Console, message string) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeProgressing, Status: operatorsv1.ConditionTrue, Reason: reasonSyncError, Message: message, LastTransitionTime: metav1.Now()})
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeDegraded, Status: operatorsv1.ConditionTrue, Message: message, Reason: reasonSyncError, LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionsManagementStateUnmanaged(operatorConfig *operatorsv1.Console) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeAvailable, Status: operatorsv1.ConditionTrue, Reason: reasonUnmanaged, Message: "The operator is in an unmanaged state, therefore its availability is unknown.", LastTransitionTime: metav1.Now()})
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeProgressing, Status: operatorsv1.ConditionFalse, Reason: reasonUnmanaged, Message: "The operator is in an unmanaged state, therefore no changes are being applied.", LastTransitionTime: metav1.Now()})
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeDegraded, Status: operatorsv1.ConditionFalse, Reason: reasonUnmanaged, Message: "The operator is in an unmanaged state, therefore no operator actions are degraded.", LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionsManagementStateRemoved(operatorConfig *operatorsv1.Console) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeAvailable, Status: operatorsv1.ConditionTrue, Reason: reasonRemoved, Message: "The operator is in a removed state, the console has been removed.", LastTransitionTime: metav1.Now()})
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeProgressing, Status: operatorsv1.ConditionFalse, Reason: reasonRemoved, Message: "The operator is in a removed state, therefore no changes are being applied.", LastTransitionTime: metav1.Now()})
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeDegraded, Status: operatorsv1.ConditionFalse, Reason: reasonRemoved, Message: "The operator is in a removed state, therefore no operator actions are degraded.", LastTransitionTime: metav1.Now()})
	return operatorConfig
}
func (c *consoleOperator) ConditionsManagementStateInvalid(operatorConfig *operatorsv1.Console) *operatorsv1.Console {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeAvailable, Status: operatorsv1.ConditionUnknown, Reason: reasonInvalid, Message: "The operator management state is invalid", LastTransitionTime: metav1.Now()})
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeProgressing, Status: operatorsv1.ConditionFalse, Reason: reasonInvalid, Message: "The operator management state is invalid", LastTransitionTime: metav1.Now()})
	v1helpers.SetOperatorCondition(&operatorConfig.Status.Conditions, operatorsv1.OperatorCondition{Type: operatorsv1.OperatorStatusTypeDegraded, Status: operatorsv1.ConditionUnknown, Reason: reasonInvalid, Message: "The operator management state is invalid", LastTransitionTime: metav1.Now()})
	return operatorConfig
}
