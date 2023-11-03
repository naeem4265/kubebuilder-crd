package apps

import (
	crdappsv1 "github.com/naeem4265/kubebuilder-crd/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// customService creates a new service for a custom resource
func customService(resource *crdappsv1.Book) corev1.Service {
	labels := map[string]string{
		"app":        trimTheOwnerPartFromImageName(resource.Spec.Container.Image),
		"controller": resource.Name,
	}
	service := corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.Spec.Service.Name,
			Namespace: resource.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(resource, crdappsv1.GroupVersion.WithKind("Book")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       resource.Spec.Container.Port,
					TargetPort: intstr.FromInt32(resource.Spec.Container.Port),
				},
			},
			Type: getTheServiceType(resource.Spec.Service.ServiceType),
		},
	}
	return service
}

func getTheServiceType(s string) corev1.ServiceType {
	if s == "NodePort" {
		return corev1.ServiceTypeNodePort
	} else if s == "ClusterIP" {
		return corev1.ServiceTypeClusterIP
	}
	return corev1.ServiceTypeClusterIP
}
