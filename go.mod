module github.com/Markus-Schwer/eventing-amqp

go 1.16

require (
	github.com/Azure/go-amqp v0.13.6
	github.com/cloudevents/sdk-go/v2 v2.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	go.uber.org/zap v1.16.0
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.19.7
	knative.dev/eventing v0.21.1-0.20210309225325-879407f613a0
	knative.dev/hack v0.0.0-20210309141825-9b73a256fd9a
	knative.dev/pkg v0.0.0-20210310050525-cc278e1666ca
)
