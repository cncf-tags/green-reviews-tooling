package main

import (
	"log"
)

func main() {
	// Step 1: Collect metrics
	metricsData, err := CollectMetrics()
	if err != nil {
		log.Fatalf("Error collecting metrics: %v", err)
	}

	// Step 2: Compute SCI
	sciValue, err := ComputeSCI(metricsData)
	if err != nil {
		log.Fatalf("Error computing SCI: %v", err)
	}

	// Step 3: Store results
	err = StoreResults(sciValue)
	if err != nil {
		log.Fatalf("Error storing results: %v", err)
	}

	log.Println("Sustainability metrics reporting completed successfully.")
}
