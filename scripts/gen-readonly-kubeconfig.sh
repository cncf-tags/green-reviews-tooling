#!/bin/bash

cat <<EOF | kubectl apply -f -
---
####################################################################
# Add custom resources which are not covered by view clusterrole
####################################################################
#
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: crole-customresources-readyonly
  labels:
    rbac.authorization.k8s.io/aggregate-to-view: "true"
rules: []
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: readonly-account
  namespace: kube-system
---

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: crolebinding-pod-readyonly
subjects:
- kind: ServiceAccount
  name: readonly-account
  namespace: kube-system
roleRef:
  kind: ClusterRole
  name: view
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: Secret
type: kubernetes.io/service-account-token
metadata:
  name: readonly-account-token
  namespace: kube-system
  annotations:
    kubernetes.io/service-account.name: readonly-account
...
EOF


USER_TOKEN_VALUE=$(kubectl -n kube-system get secret/readonly-account-token -o=go-template='{{.data.token}}' | base64 --decode)
CURRENT_CONTEXT=$(kubectl config current-context)
CURRENT_CLUSTER=$(kubectl config view --raw -o=go-template='{{range .contexts}}{{if eq .name "'''${CURRENT_CONTEXT}'''"}}{{ index .context "cluster" }}{{end}}{{end}}')
CLUSTER_CA=$(kubectl config view --raw -o=go-template='{{range .clusters}}{{if eq .name "'''${CURRENT_CLUSTER}'''"}}"{{with index .cluster "certificate-authority-data" }}{{.}}{{end}}"{{ end }}{{ end }}')
CLUSTER_SERVER=$(kubectl config view --raw -o=go-template='{{range .clusters}}{{if eq .name "'''${CURRENT_CLUSTER}'''"}}{{ .cluster.server }}{{end}}{{ end }}')

cat << EOF > green-reviews-cluster-readonly-config
apiVersion: v1
kind: Config
current-context: ${CURRENT_CONTEXT}
contexts:
- name: ${CURRENT_CONTEXT}
  context:
    cluster: ${CURRENT_CONTEXT}
    user: readonly-account
clusters:
- name: ${CURRENT_CONTEXT}
  cluster:
    certificate-authority-data: ${CLUSTER_CA}
    server: ${CLUSTER_SERVER}
users:
- name: readonly-account
  user:
    token: ${USER_TOKEN_VALUE}
EOF
