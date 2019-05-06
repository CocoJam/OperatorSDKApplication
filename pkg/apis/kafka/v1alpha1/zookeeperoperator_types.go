package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ZooKeeperOperatorSpec defines the desired state of ZooKeeperOperator
// +k8s:openapi-gen=true
type ZooKeeperOperatorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas int `json: replicas`
	ContainerName string `json: conatinerName`
	Image string `json: image`
	ContainerPorts map[string]string `json: containerPorts`
	Heap string `json: heap`
	LogDir string `json: logDir`
	DataLogDir string `json: dataLogDir`
	MountNum int `json: mountNum`
	Commands []string `json: commands`
	Args []string `json: args`
	ClientPort string `json: clientPort`
	ServerPort string `json: serverPort`
	LeaderElectionPort string `json: leaderElectionPort`
}

// ZooKeeperOperatorStatus defines the observed state of ZooKeeperOperator
// +k8s:openapi-gen=true
type ZooKeeperOperatorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZooKeeperOperator is the Schema for the zookeeperoperators API
// +k8s:openapi-gen=true
type ZooKeeperOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ZooKeeperOperatorSpec   `json:"spec,omitempty"`
	Status ZooKeeperOperatorStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ZooKeeperOperatorList contains a list of ZooKeeperOperator
type ZooKeeperOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ZooKeeperOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ZooKeeperOperator{}, &ZooKeeperOperatorList{})
}
