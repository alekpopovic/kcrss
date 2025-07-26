package main

import (
	"fmt"
	"os"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	var metricsAddr string
	var enableLeaderElection bool

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "secret-injector-controller",
	})
	if err != nil {
		fmt.Printf("unable to start manager: %v\n", err)
		os.Exit(1)
	}

	if err = (&SecretInjectorReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("SecretInjector"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		fmt.Printf("unable to create controller: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fmt.Printf("problem running manager: %v\n", err)
		os.Exit(1)
	}
}
