# Cluster Access

## Read only kubeconfig

> Refer: [Link](https://codeforphilly.github.io/chime/operations/limited-kubeconfigs/limited-kubeconfigs.html)

[create-token](https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/#create-token)

## Steps to get the kubeconfig

```bash
chmod u+x scripts/gen-readonly-kubeconfig.sh
./scripts/gen-readonly-kubeconfig.sh
```

## Test out the kubeconfig
```bash
export KUBECONFIG=${PWD}/green-reviews-cluster-readonly-config
```

```bash
kubectl get no # it will fail
```

> Just get, watch, list for all pods
