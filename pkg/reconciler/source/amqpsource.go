package source

import (
	"context"
	"fmt"

	//k8s.io imports
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1listers "k8s.io/client-go/listers/apps/v1"

	sourcesv1alpha1 "github.com/Markus-Schwer/eventing-amqp/pkg/apis/sources/v1alpha1"
	"github.com/Markus-Schwer/eventing-amqp/pkg/client/clientset/versioned"
	amqpreconciler "github.com/Markus-Schwer/eventing-amqp/pkg/client/injection/reconciler/sources/v1alpha1/amqpsource"
	listers "github.com/Markus-Schwer/eventing-amqp/pkg/client/listers/sources/v1alpha1"
	"github.com/Markus-Schwer/eventing-amqp/pkg/reconciler/source/resources"

	//knative.dev/pkg imports
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/resolver"
)

// Reconciler reconciles a AmqpSource object
type Reconciler struct {
	kubeClientSet kubernetes.Interface

	amqpLister       listers.AmqpSourceLister
	deploymentLister appsv1listers.DeploymentLister

	amqpClientSet versioned.Interface

	sinkResolver *resolver.URIResolver

	receiveAdapterImage string
}

var _ amqpreconciler.Interface = (*Reconciler)(nil)

// var _ amqpreconciler.Finalizer = (*Reconciler)(nil)

func newDeploymentCreated(namespace, name string) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeNormal, "DeploymentCreated", "AmqpSource created deployment: \"%s/%s\"", namespace, name)
}

func newDeploymentFailed(namespace, name string, err error) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeWarning, "DeploymentFailed", "AmqpSource failed to create deployment: \"%s/%s\", %w", namespace, name, err)
}

func (r *Reconciler) ReconcileKind(ctx context.Context, src *sourcesv1alpha1.AmqpSource) pkgreconciler.Event {
	src.Status.InitializeConditions()

	dest := src.Spec.Sink.DeepCopy()
	if dest.Ref != nil {
		// To call URIFromDestination(), dest.Ref must have a Namespace. If there is
		// no Namespace defined in dest.Ref, we will use the Namespace of the src
		// as the Namespace of dest.Ref.
		if dest.Ref.Namespace == "" {
			dest.Ref.Namespace = src.GetNamespace()
		}
	}

	uri, err := r.sinkResolver.URIFromDestinationV1(ctx, *dest, src)
	if err != nil {
		src.Status.MarkNoSink("NotFound", "%s", err)
		return err
	}

	src.Status.MarkSink(uri)

	ra, err := r.createReceiveAdapter(ctx, src, uri)
	if err != nil {
		logging.FromContext(ctx).Error("Unable to create the receive adapter", err.Error())
		return err
	}

	src.Status.MarkDeployed(ra)
	src.Status.CloudEventAttributes = r.createCloudEventAttributes(src)

	return nil
}

func (r *Reconciler) createReceiveAdapter(ctx context.Context, src *sourcesv1alpha1.AmqpSource, sinkURI *apis.URL) (*v1.Deployment, error) {
	raArgs := resources.ReceiveAdapterArgs{
		Image:   r.receiveAdapterImage,
		Source:  src,
		Labels:  resources.GetLabels(src.Name),
		SinkURI: sinkURI.String(),
	}
	expected := resources.MakeReceiveAdapter(&raArgs)

	ra, err := r.kubeClientSet.AppsV1().Deployments(src.Namespace).Get(ctx, expected.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		ra, err = r.kubeClientSet.AppsV1().Deployments(src.Namespace).Create(ctx, expected, metav1.CreateOptions{})
		if err != nil {
			return nil, newDeploymentFailed(ra.Namespace, ra.Name, err)
		}

		return ra, newDeploymentCreated(ra.Namespace, ra.Name)
	} else if err != nil {
		logging.FromContext(ctx).Error("Unable to get an existing receive adapter", err.Error())
		return nil, err
	} else if !metav1.IsControlledBy(ra, src) {
		return nil, fmt.Errorf("deployment %q is not owned by AmqpSource %q", ra.Name, src.Name)
	}

	return ra, nil
}

func (r *Reconciler) createCloudEventAttributes(src *sourcesv1alpha1.AmqpSource) []duckv1.CloudEventAttributes {
	ceAttribute := duckv1.CloudEventAttributes{
		Type:   sourcesv1alpha1.AmqpEventType,
		Source: sourcesv1alpha1.AmqpEventSource(src.Namespace, src.Name, src.Spec.Topic),
	}
	return []duckv1.CloudEventAttributes{ceAttribute}
}
