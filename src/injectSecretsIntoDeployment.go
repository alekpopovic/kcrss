package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)



func (r *SecretInjectorReconciler) injectSecretsIntoDeployment(deployment *appsv1.Deployment, secrets map[string]string) bool {
	updated := false
	
	for i := range deployment.Spec.Template.Spec.Containers {
		container := &deployment.Spec.Template.Spec.Containers[i]
		
		// Create a map of existing env vars for quick lookup
		existingEnvs := make(map[string]int)
		for j, env := range container.Env {
			existingEnvs[env.Name] = j
		}
		
		// Add or update environment variables
		for key, value := range secrets {
			envVar := corev1.EnvVar{
				Name:  key,
				Value: value,
			}
			
			if idx, exists := existingEnvs[key] {
				if container.Env[idx].Value != value {
					container.Env[idx] = envVar
					updated = true
				}
			} else {
				container.Env = append(container.Env, envVar)
				updated = true
			}
		}
	}
	
	return updated
}