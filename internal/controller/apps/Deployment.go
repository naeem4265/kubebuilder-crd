package apps

import (
	crdappsv1 "github.com/naeem4265/kubebuilder-crd/api/apps/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// newDeployment creates a new Deployment for a Book resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Book resource that 'owns' it.
func newDeployment(resource *crdappsv1.Book) appsv1.Deployment {
	labels := map[string]string{
		"app":        trimTheOwnerPartFromImageName(resource.Spec.Container.Image),
		"controller": resource.Name,
	}
	return appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.Spec.DeploymentName + "-controller",
			Namespace: resource.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: resource.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "my-app-pod",
							Image: resource.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: resource.Spec.Container.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}

func trimTheOwnerPartFromImageName(s string) string {
	arr := strings.Split(s, "/")
	if len(arr) == 1 {
		return arr[0]
	}
	return arr[1]
}
