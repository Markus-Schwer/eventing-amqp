package resources

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/kmeta"

	"github.com/Markus-Schwer/eventing-amqp/pkg/apis/sources/v1alpha1"
)

// ReceiveAdapterArgs are the arguments needed to create a AmqpSource Receive Adapter.
// Every field is required.
type ReceiveAdapterArgs struct {
	Image          string
	Source         *v1alpha1.AmqpSource
	Labels         map[string]string
	SinkURI        string
	AdditionalEnvs []corev1.EnvVar
}

// MakeReceiveAdapter generates (but does not insert into K8s) the Receive Adapter Deployment for
// amqp sources.
func MakeReceiveAdapter(args *ReceiveAdapterArgs) *v1.Deployment {
	replicas := int32(1)

	env := append([]corev1.EnvVar{
		{
			Name:  "AMQP_BROKER",
			Value: args.Source.Spec.Broker,
		},
		{
			Name:  "AMQP_TOPIC",
			Value: args.Source.Spec.Topic,
		},
		{
			Name: "AMQP_USER",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: args.Source.Spec.User.SecretKeyRef,
			},
		},
		{
			Name: "AMQP_PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: args.Source.Spec.Password.SecretKeyRef,
			},
		},
		{
			Name:  "SINK_URI",
			Value: args.SinkURI,
		},
		{
			Name:  "K_SINK",
			Value: args.SinkURI,
		},
		{
			Name:  "NAME",
			Value: args.Source.Name,
		},
		{
			Name:  "NAMESPACE",
			Value: args.Source.Namespace,
		},
	}, args.AdditionalEnvs...)

	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:         kmeta.ChildName(fmt.Sprintf("amqpsource-%s-", args.Source.Name), string(args.Source.UID)),
			Namespace:    args.Source.Namespace,
			GenerateName: fmt.Sprintf("%s-", args.Source.Name),
			Labels:       args.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(args.Source),
			},
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: args.Labels,
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"sidecar.istio.io/inject": "true",
					},
					Labels: args.Labels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: args.Source.Spec.ServiceAccountName,
					Containers: []corev1.Container{
						{
							Name:            "receive-adapter",
							Image:           args.Image,
							ImagePullPolicy: "IfNotPresent",
							Env:             env,
						},
					},
				},
			},
		},
	}
}
