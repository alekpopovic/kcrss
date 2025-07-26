package main

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *SecretInjectorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("secretinjector", req.NamespacedName)

	// Fetch the SecretInjector instance
	var secretInjector SecretInjector
	if err := r.Get(ctx, req.NamespacedName, &secretInjector); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Fetch secrets from remote server
	secrets, err := r.fetchSecretsFromRemote(ctx, &secretInjector)
	if err != nil {
		log.Error(err, "Failed to fetch secrets from remote server")
		r.updateStatus(ctx, &secretInjector, false, fmt.Sprintf("Failed to fetch secrets: %v", err))
		return ctrl.Result{RequeueAfter: time.Minute * 5}, nil
	}

	// Find target deployment
	var deployment appsv1.Deployment
	deploymentKey := types.NamespacedName{
		Name:      secretInjector.Spec.TargetDeployment,
		Namespace: secretInjector.Namespace,
	}

	if err := r.Get(ctx, deploymentKey, &deployment); err != nil {
		log.Error(err, "Failed to get target deployment")
		r.updateStatus(ctx, &secretInjector, false, fmt.Sprintf("Target deployment not found: %v", err))
		return ctrl.Result{RequeueAfter: time.Minute * 2}, nil
	}

	// Update deployment with secrets as env vars
	updated := r.injectSecretsIntoDeployment(&deployment, secrets)
	if updated {
		if err := r.Update(ctx, &deployment); err != nil {
			log.Error(err, "Failed to update deployment")
			r.updateStatus(ctx, &secretInjector, false, fmt.Sprintf("Failed to update deployment: %v", err))
			return ctrl.Result{RequeueAfter: time.Minute * 2}, nil
		}
		log.Info("Successfully updated deployment with secrets")
	}

	// Update status
	r.updateStatus(ctx, &secretInjector, true, "Secrets successfully injected")

	// Schedule next reconciliation
	interval := time.Minute * 10 // default
	if secretInjector.Spec.RefreshInterval != "" {
		if parsedInterval, err := time.ParseDuration(secretInjector.Spec.RefreshInterval); err == nil {
			interval = parsedInterval
		}
	}

	return ctrl.Result{RequeueAfter: interval}, nil
}
