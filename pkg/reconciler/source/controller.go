package source

import (
	"context"

	"github.com/kelseyhightower/envconfig"

	//k8s.io imports
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"

	//Injection imports
	sourcescheme "github.com/Markus-Schwer/eventing-amqp/pkg/client/clientset/versioned/scheme"
	amqpclient "github.com/Markus-Schwer/eventing-amqp/pkg/client/injection/client"
	amqpinformer "github.com/Markus-Schwer/eventing-amqp/pkg/client/injection/informers/sources/v1alpha1/amqpsource"
	amqpreconciler "github.com/Markus-Schwer/eventing-amqp/pkg/client/injection/reconciler/sources/v1alpha1/amqpsource"
	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"

	//knative.dev/eventing imports
	"knative.dev/eventing/pkg/apis/sources/v1alpha1"

	//knative.dev/pkg imports
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/resolver"
)

// envConfig will be used to extract the required environment variables.
// If this configuration cannot be extracted, then NewController will panic.
type envConfig struct {
	Image string `envconfig:"AMQP_RA_IMAGE"`
}

func NewController(ctx context.Context, _ configmap.Watcher) *controller.Impl {
	env := &envConfig{}
	if err := envconfig.Process("", env); err != nil {
		logging.FromContext(ctx).Panicf("unable to process AmqpSource's required environment variables: %v", err)
	}

	amqpInformer := amqpinformer.Get(ctx)
	deploymentInformer := deploymentinformer.Get(ctx)

	r := &Reconciler{
		kubeClientSet:       kubeclient.Get(ctx),
		amqpClientSet:       amqpclient.Get(ctx),
		deploymentLister:    deploymentInformer.Lister(),
		receiveAdapterImage: env.Image, // can be empty
	}
	impl := amqpreconciler.NewImpl(ctx, r)

	r.sinkResolver = resolver.NewURIResolver(ctx, impl.EnqueueKey)

	logging.FromContext(ctx).Info("Setting up AMQP event handlers")

	// Watch for changes from any AmqpSource object
	amqpInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	deploymentInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterControllerGK(v1alpha1.Kind("AmqpSource")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}

func init() {
	sourcescheme.AddToScheme(scheme.Scheme)
}
