package manifest

const (
	ClusterlinkNetworkManagerClusterRole = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: clusterlink-network-manager
rules:
  - apiGroups: ['*']
    resources: ['*']
    verbs: ["*"]
  - nonResourceURLs: ['*']
    verbs: ["get"]
`

	ClusterlinkFloaterClusterRole = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: clusterlink-floater
rules:
  - apiGroups: ['*']
    resources: ['*']
    verbs: ["*"]
  - nonResourceURLs: ['*']
    verbs: ["get"]
`

	KosmosClusterRole = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kosmos
rules:
  - apiGroups: ['*']
    resources: ['*']
    verbs: ["*"]
  - nonResourceURLs: ['*']
    verbs: ["*"]
`

	ClusterTreeClusterRole = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: clustertree
rules:
  - apiGroups: ['*']
    resources: ['*']
    verbs: ["*"]
  - nonResourceURLs: ['*']
    verbs: ["get"]
`

	CorednsClusterRole = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kosmos-coredns
rules:
  - apiGroups: ['*']
    resources: ['*']
    verbs: ["*"]
  - nonResourceURLs: ['*']
    verbs: ["get"]
`

	SchedulerClusterRole = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kosmos-scheduler
rules:
  - verbs:
      - get
      - list
      - watch
    apiGroups:
      - kosmos.io
    resources:
      - '*'
  - verbs:
      - create
      - patch
      - update
    apiGroups:
      - ''
      - events.k8s.io
    resources:
      - events
  - verbs:
      - create
    apiGroups:
      - coordination.k8s.io
    resources:
      - leases
  - verbs:
      - get
      - update
    apiGroups:
      - coordination.k8s.io
    resources:
      - leases
    resourceNames:
      - kosmos-scheduler
  - verbs:
      - create
    apiGroups:
      - ''
    resources:
      - endpoints
  - verbs:
      - get
      - update
    apiGroups:
      - ''
    resources:
      - endpoints
  - verbs:
      - get
      - list
      - watch
    apiGroups:
      - ''
    resources:
      - nodes
  - verbs:
      - delete
      - get
      - list
      - watch
    apiGroups:
      - ''
    resources:
      - pods
  - verbs:
      - create
    apiGroups:
      - ''
    resources:
      - bindings
      - pods/binding
  - verbs:
      - patch
      - update
    apiGroups:
      - ''
    resources:
      - pods/status
  - verbs:
      - get
      - list
      - watch
    apiGroups:
      - ''
    resources:
      - replicationcontrollers
      - services
  - verbs:
      - get
      - list
      - watch
    apiGroups:
      - apps
      - extensions
    resources:
      - replicasets
  - verbs:
      - get
      - list
      - watch
    apiGroups:
      - apps
    resources:
      - statefulsets
  - verbs:
      - get
      - list
      - watch
    apiGroups:
      - policy
    resources:
      - poddisruptionbudgets
  - verbs:
      - get
      - list
      - watch
      - update
    apiGroups:
      - ''
    resources:
      - persistentvolumeclaims
      - persistentvolumes
  - verbs:
      - create
    apiGroups:
      - authentication.k8s.io
    resources:
      - tokenreviews
  - verbs:
      - create
    apiGroups:
      - authorization.k8s.io
    resources:
      - subjectaccessreviews
  - verbs:
      - get
      - list
      - watch
    apiGroups:
      - storage.k8s.io
    resources:
      - '*'
  - verbs:
      - get
      - list
      - watch
    apiGroups:
      - ''
    resources:
      - configmaps
      - namespaces
`
)
