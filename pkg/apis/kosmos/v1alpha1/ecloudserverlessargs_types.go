package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:defaulter-gen=true

// ECloudServerlessArgs holds arguments used to configure the ECloudServerless node provider
type ECloudServerlessArgs struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	AccessKey string `json:"accessKey,omitempty"`

	// +optional
	SecretKey string `json:"secretKey,omitempty"`
}
