package cmd

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

func WaitForNamespace(namespace string) []string {
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
