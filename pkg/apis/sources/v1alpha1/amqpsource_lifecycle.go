package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	"knative.dev/eventing/pkg/apis/duck"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

const (
	AmqpConditionReady                           = apis.ConditionReady
	AmqpConditionSinkProvided apis.ConditionType = "SinkProvided"
	AmqpConditionDeployed     apis.ConditionType = "Deployed"
	AmqpConditionResources    apis.ConditionType = "ResourcesReady"
)

var AmqpSourceCondSet = apis.NewLivingConditionSet(AmqpConditionSinkProvided, AmqpConditionDeployed)

func (s *AmqpSourceStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return AmqpSourceCondSet.Manage(s).GetCondition(t)
}

func (s *AmqpSourceStatus) GetTopLevelCondition() *apis.Condition {
	return AmqpSourceCondSet.Manage(s).GetTopLevelCondition()
}

func (in *AmqpSource) GetStatus() *duckv1.Status {
	return &in.Status.Status
}

func (in *AmqpSource) GetConditionSet() apis.ConditionSet {
	return AmqpSourceCondSet
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (s *AmqpSourceStatus) InitializeConditions() {
	AmqpSourceCondSet.Manage(s).InitializeConditions()
}

// IsReady returns if the source is ready
func (s *AmqpSourceStatus) IsReady() bool {
	return AmqpSourceCondSet.Manage(s).IsHappy()
}

// MarkSink sets the condition that the source has a sink configured.
func (s *AmqpSourceStatus) MarkSink(uri *apis.URL) {
	s.SinkURI = uri
	if len(uri.String()) > 0 {
		AmqpSourceCondSet.Manage(s).MarkTrue(AmqpConditionSinkProvided)
	} else {
		AmqpSourceCondSet.Manage(s).MarkUnknown(AmqpConditionSinkProvided, "SinkEmpty", "Sink has resolved to empty.%s", "")
	}
}

// MarkNoSink sets the condition that the source does not have a sink configured.
func (s *AmqpSourceStatus) MarkNoSink(reason, messageFormat string, messageA ...interface{}) {
	AmqpSourceCondSet.Manage(s).MarkFalse(AmqpConditionSinkProvided, reason, messageFormat, messageA...)
}

func DeploymentIsAvailable(d *appsv1.DeploymentStatus, def bool) bool {
	for _, cond := range d.Conditions {
		if cond.Type == appsv1.DeploymentAvailable {
			return cond.Status == "True"
		}
	}
	return def
}

// MarkDeployed marks the source's Deployed condition to True
func (s *AmqpSourceStatus) MarkDeployed(d *appsv1.Deployment) {
	if duck.DeploymentIsAvailable(&d.Status, false) {
		AmqpSourceCondSet.Manage(s).MarkTrue(AmqpConditionDeployed)
	} else {
		AmqpSourceCondSet.Manage(s).MarkFalse(AmqpConditionDeployed, "DeploymentUnavailable", "The Deployment '%s' is unavailable.", d.Name)
	}
}

// MarkNotDeployed marks the source's Deployed condition to Deploying with
// the provided reason and message.
func (s *AmqpSourceStatus) MarkDeploying(reason, messageFormat string, messageA ...interface{}) {
	AmqpSourceCondSet.Manage(s).MarkUnknown(AmqpConditionDeployed, reason, messageFormat, messageA...)
}

// MarkNotDeployed marks the source's Deployed condition to False with
// the provided reason and message.
func (s *AmqpSourceStatus) MarkNotDeployed(reason, messageFormat string, messageA ...interface{}) {
	AmqpSourceCondSet.Manage(s).MarkFalse(AmqpConditionDeployed, reason, messageFormat, messageA...)
}

// MarkResourcesCorrect marks the source's Resources condition to Correct
func (s *AmqpSourceStatus) MarkResourcesCorrect() {
	AmqpSourceCondSet.Manage(s).MarkTrue(AmqpConditionResources)
}

// MarkResourcesIncorrect marks the source's Resources condition to Incorrect with
// the provided reason and message.
func (s *AmqpSourceStatus) MarkResourcesIncorrect(reason, messageFormat string, messageA ...interface{}) {
	AmqpSourceCondSet.Manage(s).MarkFalse(AmqpConditionResources, reason, messageFormat, messageA...)
}
