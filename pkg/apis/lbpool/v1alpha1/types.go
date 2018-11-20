package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type LoadBalancer struct {
	metav1.TypeMeta   `json: ", inline"`
	metav1.ObjectMeta `json: "metadata, omitempty"`

	Spec   LoadBalancerSpec   `json: "spec"`
	Status LoadBalancerStatus `json: "status"`
}

type LoadBalancerSpec struct {
	DeploymentName string `json:"deploymentName"`
	Replicas       *int32 `json:"replicas"`
}

type LoadBalancerStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type LoadBalancerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []LoadBalancer `json: "items"`
}
