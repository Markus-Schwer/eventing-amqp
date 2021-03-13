package main

import (
	amqpadapter "github.com/Markus-Schwer/eventing-amqp/pkg/adapter"

	"knative.dev/eventing/pkg/adapter/v2"
)

func main() {
	adapter.Main("amqpsource", amqpadapter.NewEnvConfig, amqpadapter.NewAdapter)
}
