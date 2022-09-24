package controllers

import (
	"context"

	"github.com/dkhachyan/jupyterhub-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *JupyterhubReconciler) jupyerhubService(cr *v1alpha1.Jupyterhub) *corev1.Service {
	labels := map[string]string{
		"app": "jupyterhub",
	}
	ports := []corev1.ServicePort{
		{
			Name:       "jupyterhub-http",
			Protocol:   corev1.Protocol("TCP"),
			Port:       8000,
			TargetPort: intstr.FromInt(8000),
		},
	}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: ports,
		},
	}
	return svc
}

func (r *JupyterhubReconciler) ensureJupyterhubService(instance *v1alpha1.Jupyterhub) (*reconcile.Result, error) {
	desired := r.jupyerhubService(instance)
	existing := &corev1.Service{}

	err := r.Create(context.TODO(), desired)
	if err != nil && errors.IsAlreadyExists(err) {
		err := r.Client.Get(context.TODO(), client.ObjectKeyFromObject(desired), existing)
		if err != nil {
			return nil, err
		}
		if !apiequality.Semantic.DeepEqual(existing, desired) {
			existing.Spec = desired.Spec
			existing.Labels = desired.Labels
			err := r.Update(context.TODO(), existing)
			return nil, err
		}
	}
	return nil, nil
}
