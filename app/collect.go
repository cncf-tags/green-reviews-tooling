package main

type Metric struct {
	Name  string
	Value float64
}

func CollectMetrics() ([]Metric, error) {
	// Fetch metrics from Promehteus:
	// container_cpu_usage_seconds_total
	// container_memory_rss
	// container_memory_working_set_bytes
	// kepler_container_joules_total
	metrics := []Metric{
		// {"EnergyConsumption", energy},
		// {"CarbonIntensity", ci},
		// other metrics from prometheus
	}
	return metrics, nil
}
