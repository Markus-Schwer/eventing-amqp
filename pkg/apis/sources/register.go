package sources

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/apis/duck"
)

const (
	GroupName = "sources.knative.dev"

	// SourceDuckLabelKey is the label key to indicate
	// whether the CRD is a Source duck type.
	// Valid values: "true" or "false"
	SourceDuckLabelKey = duck.GroupName + "/source"

	// SourceDuckLabelValue is the label value to indicate
	// the CRD is a Source duck type.
	SourceDuckLabelValue = "true"
)

var (
	// ContainerSourceResource respresents a Knative Eventing Sources ContainerSource
	AmqpResource = schema.GroupResource{
		Group:    GroupName,
		Resource: "amqpsources",
	}
)
