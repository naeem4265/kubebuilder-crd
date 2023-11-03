/*
Copyright 2023.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type ContainerSpec struct {
	Image string `json:"image,omitempty"`
	Port  int32  `json:"port,omitempty"`
}

// BookSpec defines the desired state of Book
type BookSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	DeploymentName string `json:"deploymentName"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// DeploymentName represents the name of the deployment we will create using CustomCrd
	// Replicas defines number of pods will be running in the deployment
	Replicas *int32 `json:"replicas"`
	// Container contains Image and Port
	Container ContainerSpec `json:"container"`
	// Service contains ServiceName, ServiceType, ServiceNodePort
	// +optional
	Service ServiceSpec `json:"service,omitempty"`
}

// BookStatus defines the observed state of Book
type BookStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
	ServiceCreated    bool  `json:"serviceCreated"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Book is the Schema for the books API
type Book struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BookSpec   `json:"spec,omitempty"`
	Status BookStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BookList contains a list of Book
type BookList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Book `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Book{}, &BookList{})
}

// ServiceSpec defines the desired state of Service
type ServiceSpec struct {
	Name string `json:"name"`
	// +optional
	ServiceType string `json:"serviceType,omitempty"`
	ServicePort int32  `json:"servicePort"`
	// +optional
	ServiceNodePort int32 `json:"serviceNodePort,omitempty"`
}
