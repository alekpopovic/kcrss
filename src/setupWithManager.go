package main

import (
	appsv1 "k8s.io/api/apps/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func (r *SecretInjectorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&SecretInjector{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
