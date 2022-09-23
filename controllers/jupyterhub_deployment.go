package controllers

import (
	"context"

	"github.com/dkhachyan/jupyterhub-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *JupyterhubReconciler) jupyerhubDeployment(cr *v1alpha1.Jupyterhub) *appsv1.Deployment {
	labels := map[string]string{
		"app": "jupyterhub",
	}
	jupyterhubDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
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
	desired := r.jupyerhubDeployment(instance)
	existing := &appsv1.Deployment{}

	err := r.Create(context.TODO(), desired)
	if err != nil && errors.IsAlreadyExists(err) {
		err := r.Client.Get(context.TODO(), client.ObjectKeyFromObject(desired), existing)
		if err != nil {
			return nil, err
		}
		if !apiequality.Semantic.DeepEqual(existing, desired) {
			existing.Spec.Replicas = desired.Spec.Replicas
			existing.Spec.Template = desired.Spec.Template
			existing.Labels = desired.Labels
			err := r.Update(context.TODO(), existing)
			return nil, err
		}
	}
	return nil, nil
}
