package main

import (
	"knative.dev/eventing/pkg/adapter/v2"

	amqpadapter "github.com/Markus-Schwer/eventing-amqp/pkg/adapter"
)

func main() {
	adapter.MainMessageAdapter("amqpsource", amqpadapter.NewEnvConfig, amqpadapter.NewAdapter)
}
