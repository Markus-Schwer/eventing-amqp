package main

import (
	amqp "github.com/Markus-Schwer/eventing-amqp/pkg/reconciler/source"

	"knative.dev/pkg/injection/sharedmain"
)

func main() {
	sharedmain.Main("amqp-controller", amqp.NewController)
}
