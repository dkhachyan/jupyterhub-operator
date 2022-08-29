package controllers

import (
	"context"

	"github.com/dkhachyan/jupyterhub-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *JupyterhubReconciler) jupyerhubDeployment(cr *v1alpha1.Jupyterhub) *appsv1.Deployment {
	jupyterhubDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &cr.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "jupyterhub",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "jupyterhub",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "jupyterhub",
							Image: cr.Spec.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8000,
								},
							},
						},
					},
				},
			},
		},
	}
	return jupyterhubDeployment
}

func (r *JupyterhubReconciler) ensureJupyterhubDeployment(instance *v1alpha1.Jupyterhub) (*reconcile.Result, error) {
	jupyterhubDeployment := r.jupyerhubDeployment(instance)
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      jupyterhubDeployment.Name,
		Namespace: instance.Namespace,
	}, jupyterhubDeployment)

	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		err = r.Create(context.TODO(), jupyterhubDeployment)

		if err != nil {
			// Deployment failed
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		return &reconcile.Result{}, err
	}
	return nil, nil
}
