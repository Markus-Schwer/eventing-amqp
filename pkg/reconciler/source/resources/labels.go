package resources

const (
	controllerAgentName = "amqp-source-controller"
)

func GetLabels(name string) map[string]string {
	return map[string]string{
		"eventing.knative.dev/source":     controllerAgentName,
		"eventing.knative.dev/SourceName": name,
	}
}
