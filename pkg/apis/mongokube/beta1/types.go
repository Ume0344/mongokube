package beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Mk struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MkSpec   `json:"spec"`
	Status MkStatus `json:"status"`
}

type MkSpec struct {
	MongoExpressImage       string `json:"mongoExpressImage"`
	MongoExpressServicePort string `json:"mongoExpressServicePort"`
	MongoDbImage            string `json:"mongoDbImage"`
	DbUsername              string `json:"dbUsername"`
	DbPassword              string `json:"dbPassword"`
}

type MkStatus struct {
	Progress string `json:"progress"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Mk `json:"items"`
}
