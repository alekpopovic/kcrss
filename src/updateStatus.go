package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *SecretInjectorReconciler) updateStatus(ctx context.Context, si *SecretInjector, ready bool, message string) {
	si.Status.LastSync = metav1.Now()
	si.Status.Ready = ready
	si.Status.Message = message
	r.Status().Update(ctx, si)
}
