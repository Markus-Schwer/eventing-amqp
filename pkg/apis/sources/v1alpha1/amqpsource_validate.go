package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
	"knative.dev/pkg/kmp"
)

func (current *AmqpSource) Validate(ctx context.Context) *apis.FieldError {
	if apis.IsInUpdate(ctx) {
		original := apis.GetBaseline(ctx).(*AmqpSource)
		if diff, err := kmp.ShortDiff(original.Spec, current.Spec); err != nil {
			return &apis.FieldError{
				Message: "Failed to diff AmqpSource",
				Paths:   []string{"spec"},
				Details: err.Error(),
			}
		} else if diff != "" {
			return &apis.FieldError{
				Message: "Immutable fields changed (-old +new)",
				Paths:   []string{"spec"},
				Details: diff,
			}
		}
	}

	return nil
}
