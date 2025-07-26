package main

import (
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// SecretInjectorSpec defines the desired state
type SecretInjectorSpec struct {
	RemoteServerURL  string            `json:"remoteServerURL"`
	TargetDeployment string            `json:"targetDeployment"`
	SecretKeys       []string          `json:"secretKeys"`
	Headers          map[string]string `json:"headers,omitempty"`
	RefreshInterval  string            `json:"refreshInterval,omitempty"`
}

// SecretInjectorStatus defines the observed state
type SecretInjectorStatus struct {
	LastSync   metav1.Time `json:"lastSync,omitempty"`
	SecretHash string      `json:"secretHash,omitempty"`
	Ready      bool        `json:"ready"`
	Message    string      `json:"message,omitempty"`
}

// SecretInjector is the Schema for the secretinjectors API
type SecretInjector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SecretInjectorSpec   `json:"spec,omitempty"`
	Status            SecretInjectorStatus `json:"status,omitempty"`
}

// SecretInjectorList contains a list of SecretInjector
type SecretInjectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretInjector `json:"items"`
}

// SecretInjectorReconciler reconciles a SecretInjector object
type SecretInjectorReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}
