package main

func ComputeSCI(metrics []Metric) (float64, error) {
	var total float64
	for _, metric := range metrics {
		total += metric.Value
	}
	return total / float64(len(metrics)), nil
}
