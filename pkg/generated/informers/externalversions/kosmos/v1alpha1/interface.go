// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/kosmos.io/kosmos/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// Clusters returns a ClusterInformer.
	Clusters() ClusterInformer
	// ClusterNodes returns a ClusterNodeInformer.
	ClusterNodes() ClusterNodeInformer
	// DaemonSets returns a DaemonSetInformer.
	DaemonSets() DaemonSetInformer
	// Knodes returns a KnodeInformer.
	Knodes() KnodeInformer
	// NodeConfigs returns a NodeConfigInformer.
	NodeConfigs() NodeConfigInformer
	// PodConvertPolicies returns a PodConvertPolicyInformer.
	PodConvertPolicies() PodConvertPolicyInformer
	// ShadowDaemonSets returns a ShadowDaemonSetInformer.
	ShadowDaemonSets() ShadowDaemonSetInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// Clusters returns a ClusterInformer.
func (v *version) Clusters() ClusterInformer {
	return &clusterInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// ClusterNodes returns a ClusterNodeInformer.
func (v *version) ClusterNodes() ClusterNodeInformer {
	return &clusterNodeInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// DaemonSets returns a DaemonSetInformer.
func (v *version) DaemonSets() DaemonSetInformer {
	return &daemonSetInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// Knodes returns a KnodeInformer.
func (v *version) Knodes() KnodeInformer {
	return &knodeInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// NodeConfigs returns a NodeConfigInformer.
func (v *version) NodeConfigs() NodeConfigInformer {
	return &nodeConfigInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// PodConvertPolicies returns a PodConvertPolicyInformer.
func (v *version) PodConvertPolicies() PodConvertPolicyInformer {
	return &podConvertPolicyInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// ShadowDaemonSets returns a ShadowDaemonSetInformer.
func (v *version) ShadowDaemonSets() ShadowDaemonSetInformer {
	return &shadowDaemonSetInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
