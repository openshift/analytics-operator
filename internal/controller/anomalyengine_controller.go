/*
Copyright 2024 Redhat.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	v1alpha1 "github.com/openshift/analytics-operator/api/v1alpha1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// AnomalyEngineReconciler reconciles a AnomalyEngine object
type AnomalyEngineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// common to all components deployed by operator
//+kubebuilder:rbac:groups=core,resources=namespaces;services;configmaps;secrets;serviceaccounts,verbs=list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=*
//+kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=*,verbs=*

// RBAC for running anomaly operator
//+kubebuilder:rbac:groups=observability-analytics.redhat.com,resources=anomalyengines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=observability-analytics.redhat.com,resources=anomalyengines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=observability-analytics.redhat.com,resources=anomalyengines/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *AnomalyEngineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	_ = context.Background()

	var log = logf.Log.WithName("Engine Reconciler")
	log.Info("Reconcile called")

	anomaly := &v1alpha1.AnomalyEngine{}

	err := r.Get(ctx, req.NamespacedName, anomaly)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Anomaly resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed")
		return ctrl.Result{}, err
	}

	log.Info("Anomaly data", "Spec", anomaly.Spec)

	// Namespace
	var result *ctrl.Result
	result, err = r.ensureNamespace(req, anomaly)
	if result != nil {
		return *result, err
	}

	//  Storage Admin Role
	result, err = r.ensureRole(req, anomaly, r.getStorageAdminRole(anomaly))
	if result != nil {
		return *result, err
	}

	//  Service Account for CRD Storage operation
	result, err = r.ensureServiceAccount(req, anomaly, r.getServiceAccountForStorage(anomaly), false)
	if result != nil {
		return *result, err
	}

	//  Role Binding for Storage role
	result, err = r.ensureRoleBinding(req, anomaly, r.getRoleBindingForStorageServiceAccount(anomaly))
	if result != nil {
		return *result, err
	}

	log.Info("Storage logic executed")

	// ServiceAccount for the Anomaly Engine
	result, err = r.ensureServiceAccount(
		req,
		anomaly,
		r.getServiceAccountForAnomalyEngine(anomaly),
		true,
	)
	if result != nil {
		return *result, err
	}

	// ClusterRoleBinding for the Anomaly Engine
	result, err = r.ensureClusterRoleBinding(req, anomaly, r.getRoleBindingForAnomalyEngine(anomaly))
	if result != nil {
		return *result, err
	}

	// ConfigMap for the Anomaly Engine
	result, err = r.ensureConfigMap(req, anomaly)
	if result != nil {
		return *result, err
	}

	// CronJob
	result, err = r.ensureCronJob(req, anomaly)
	if result != nil {
		return *result, err
	}

	log.Info("Anomaly Engine logic executed")

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AnomalyEngineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.AnomalyEngine{}).
		Complete(r)
}
