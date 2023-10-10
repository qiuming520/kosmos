package manifest

const (
	ClusterlinkNetworkManagerDeployment = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clusterlink-network-manager
  namespace: {{ .Namespace }}
  labels:
    app: clusterlink-network-manager
spec:
  replicas: 1
  selector:
    matchLabels:
      app: clusterlink-network-manager
  template:
    metadata:
      labels:
        app: clusterlink-network-manager
    spec:
      serviceAccountName: clusterlink-network-manager
      containers:
        - name: manager
          image: {{ .ImageRepository }}/clusterlink-network-manager:v{{ .Version }}
          imagePullPolicy: IfNotPresent
          command:
            - clusterlink-network-manager
            - v=4
          resources:
            limits:
              memory: 500Mi
              cpu: 500m
            requests:
              cpu: 500m
              memory: 500Mi
`

	ClusterlinkOperatorDeployment = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clusterlink-operator
  namespace: {{ .Namespace }}
  labels:
    app: operator
spec:
  replicas: 1
  selector:
    matchLabels:
      app: operator
  template:
    metadata:
      labels:
        app: operator
    spec:
      serviceAccountName: clusterlink-operator
      affinity:
        podAntiAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            - labelSelector:
                matchExpressions:
                  - key: app
                    operator: In
                    values:
                      - operator
              namespaces:
                - clusterlink-system
              topologyKey: kubernetes.io/hostname
      containers:
      - name: operator
        image: ghcr.io/kosmos-io/clusterlink-operator:v{{ .Version }}
        imagePullPolicy: IfNotPresent
        command:
          - clusterlink-operator
          - --controlpanelconfig=/etc/clusterlink/kubeconfig
        resources:
          limits:
            memory: 200Mi
            cpu: 250m
          requests:
            cpu: 100m
            memory: 200Mi
        env:
        - name: VERSION
          value: v{{ .Version }}
        - name: CLUSTER_NAME
          value: {{ .ClusterName }}
        - name: USE_PROXY
          value: "{{ .UseProxy }}"
        volumeMounts:
          - mountPath: /etc/clusterlink
            name: proxy-config
            readOnly: true
      volumes:
        - name: proxy-config
          secret:
            secretName: controlpanel-config

`

	ClusterTreeKnodeManagerDeployment = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clustertree-knode-manager
  namespace: {{ .Namespace }}
  labels:
    app: clustertree-knode-manager
spec:
  replicas: 1
  selector:
    matchLabels:
      app: clustertree-knode-manager
  template:
    metadata:
      labels:
        app: clustertree-knode-manager
    spec:
      serviceAccountName: clustertree-knode-manager
      containers:
        - name: manager
          image: {{ .ImageRepository }}/clustertree-knode-manager:v{{ .Version }}
          imagePullPolicy: IfNotPresent
          command:
            - clustertree-knode-manager
            - --kube-api-qps=500
            - --kube-api-burst=1000
            - --kubeconfig=/etc/kube/config
            - --leader-elect=false
          volumeMounts:
            - mountPath: /etc/kube
              name: config-volume
              readOnly: true
      volumes:
        - configMap:
            defaultMode: 420
            items:
              - key: kubeconfig
                path: config
            name: host-kubeconfig
          name: config-volume
`
)

type DeploymentReplace struct {
	Namespace       string
	ImageRepository string
	Version         string
}

type ClusterlinkDeploymentReplace struct {
	Namespace   string
	Version     string
	ClusterName string
	UseProxy    string
}