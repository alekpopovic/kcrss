package main

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func init() {
	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
	appsv1.AddToScheme(scheme)
}
