package cmd

import (
	"fmt"
	"net/url"
)

func Apply(manifest string) []string {
	return []string{
		"kubectl",
		"apply",
		"-f",
		manifest,
	}
}

func Delete(manifest string) []string {
	return []string{
		"kubectl",
		"delete",
		"-f",
		manifest,
		"--wait",
	}
}

func Echo(msg string) []string {
	return []string{
		"echo",
		"'" + msg + "'",
	}
}

func FluxInstall() []string {
	return []string{
		"flux",
		"install",
	}
}

func FluxReconcile(resource, name string) []string {
	return []string{
		"flux",
		"reconcile",
		resource,
		name,
	}
}

func GetNodeNames() []string {
	return []string{
		"kubectl",
		"get",
		"node",
		"-o",
		"name",
	}
}

func LabelNode(nodeName string, labels map[string]string) []string {
	args := []string{
		"kubectl",
		"label",
		nodeName,
	}
	for k, v := range labels {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}
	return args
}

func Patch(resource, name, namespace, path, value string) []string {
	return []string{
		"kubectl",
		"patch",
		resource,
		name,
		"-n",
		namespace,
		"--type=json",
		"-p",
		fmt.Sprintf(`[{"op": "add", "path": "%s", "value": %s}]`, path, value),
	}
}

func WaitForReadyPods(namespace string) []string {
	return []string{
		"kubectl",
		"wait",
		"pod",
		"--all",
		"--namespace",
		namespace,
		"--timeout",
		"300s",
		"--for",
		"condition=Ready",
	}
}

// QueryPrometheus executes an instant PromQL query via kubectl's API server proxy,
// returning the raw JSON response from the Prometheus HTTP API.
func QueryPrometheus(query string) []string {
	path := "/api/v1/namespaces/monitoring/services/prometheus-kube-prometheus-prometheus:9090/proxy/api/v1/query?query=" + url.QueryEscape(query)
	return []string{
		"kubectl",
		"get",
		"--raw",
		path,
	}
}
