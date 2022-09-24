/*
Copyright 2022.

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

package controllers

import (
	"context"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/dkhachyan/jupyterhub-operator/api/v1alpha1"
)

// JupyterhubReconciler reconciles a Jupyterhub object
type JupyterhubReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=jupyter.org,resources=jupyterhubs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=jupyter.org,resources=jupyterhubs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=jupyter.org,resources=jupyterhubs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Jupyterhub object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *JupyterhubReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	log := log.FromContext(ctx).WithValues("Jupyterhub", req.NamespacedName)

	instance := &v1alpha1.Jupyterhub{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	result, err := r.ensureJupyterhub(instance)
	if result != nil {
		log.Error(err, "Deployment Not ready")
		return *result, err
	}
	return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JupyterhubReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Jupyterhub{}).
		Complete(r)
}

func (r *JupyterhubReconciler) ensureJupyterhub(instance *v1alpha1.Jupyterhub) (result *ctrl.Result, err error) {

	result, err = r.ensureJupyterhubDeployment(instance)
	if err != nil {
		return result, err
	}

	result, err = r.ensureJupyterhubService(instance)
	if err != nil {
		return result, err
	}

	return result, err
}
