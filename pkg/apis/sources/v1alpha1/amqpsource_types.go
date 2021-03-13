package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AmqpSource is the Schema for the amqpsources API.
// +k8s:openapi-gen=true
type AmqpSource struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec holds the desired state of the AmqpSource (from the client).
	Spec AmqpSourceSpec `json:"spec"`

	// Status communicates the observed state of the AmqpSource (from the controller).
	// +optional
	Status AmqpSourceStatus `json:"status,omitempty"`
}

// var _ runtime.Object = (*AmqpSource)(nil)
// var _ resourcesemantics.GenericCRD = (*AmqpSource)(nil)
// var _ kmeta.OwnerRefable = (*AmqpSource)(nil)
// var _ apis.Defaultable = (*AmqpSource)(nil)
// var _ apis.Validatable = (*AmqpSource)(nil)
// var _ duckv1.KRShaped = (*AmqpSource)(nil)

// AmqpSourceSpec holds the desired state of the AmqpSource (from the client).
type AmqpSourceSpec struct {
	// Broker is the Amqp server the consumer will connect to.
	// +required
	Broker string `json:"broker"`
	// Topic topic to consume messages from
	// +required
	Topic string `json:"topic,omitempty"`
	// User for amqp connection
	// +optional
	User SecretValueFromSource `json:"user,omitempty"`
	// Password for amqp connection
	// +optional
	Password SecretValueFromSource `json:"password,omitempty"`
	// Sink is a reference to an object that will resolve to a domain name to use as the sink.
	// +optional
	Sink *duckv1.Destination `json:"sink,omitempty"`
	// ServiceAccountName is the name of the ServiceAccount that will be used to run the Receive
	// Adapter Deployment.
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
}

// SecretValueFromSource represents the source of a secret value
type SecretValueFromSource struct {
	// The Secret key to select from.
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

const (
	AmqpEventType = "dev.knative.amqp.event"
)

func AmqpEventSource(namespace, amqpSourceName, topic string) string {
	return fmt.Sprintf("/apis/v1/namespaces/%s/amqpsources/%s#%s", namespace, amqpSourceName, topic)
}

type AmqpSourceStatus struct {
	// inherits duck/v1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1.SourceStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AmqpSourceList contains a list of AmqpSources.
type AmqpSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AmqpSource `json:"items"`
}
