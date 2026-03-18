// Package sci computes the Software Carbon Intensity (SCI) score from energy
// measurements and provides the types used to represent benchmark results.
package sci

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// wattSecondsToKWh is a fixed unit conversion (Joules → kWh), not infrastructure-specific.
const wattSecondsToKWh = 0.000000277777777777778

// SCIConfig holds the infrastructure-specific inputs to the SCI formula.
// The defaults below match the Grafana dashboard in clusters/base/falco-sci.yaml
// (Equinix Paris / France grid) and are pending update for Oracle Cloud / US-Ashburn.
type SCIConfig struct {
	EmissionsFactor EmissionsComponent
	PUE             PUEComponent
	EmbodiedCarbon  EmbodiedComponent
}

// DefaultSCIConfig returns the current infrastructure defaults.
// TODO: source these from projects.json or workflow inputs once Oracle Cloud
// values (region, PUE, embodied carbon for BM.Standard2.52) are established.
func DefaultSCIConfig() SCIConfig {
	return SCIConfig{
		EmissionsFactor: EmissionsComponent{
			Region: "France",
			Units:  "gCO2eq/kWh",
			Value:  78.81, // 2023 annual average, CO2.js/Ember data
		},
		PUE: PUEComponent{
			Provider: "Equinix",
			Value:    1.42, // Equinix 2023 sustainability report
		},
		EmbodiedCarbon: EmbodiedComponent{
			Units: "gCO2eq",
			Value: 3.92, // 15-min run on Equinix m3.small.x86, Boavizta API
		},
	}
}

// BenchmarkResult holds the full results of a pipeline run per proposal-003.
type BenchmarkResult struct {
	CNCFProject ProjectInfo `json:"cncf_project"`
	Results     Results     `json:"results"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     time.Time   `json:"end_time"`
}

type ProjectInfo struct {
	Name    string `json:"name"`
	Config  string `json:"config"`
	Version string `json:"version"`
}

type Results struct {
	Metrics []Metric `json:"metrics"`
	SCI     SCI      `json:"sci"`
}

type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type SCI struct {
	Components SCIComponents `json:"components"`
	Units      string        `json:"units"`
	Value      float64       `json:"value"`
}

type SCIComponents struct {
	EnergyConsumption EnergyComponent    `json:"energy_consumption"`
	EmissionsFactor   EmissionsComponent `json:"emissions_factor"`
	EmbodiedCarbon    EmbodiedComponent  `json:"embodied_carbon"`
	PUE               PUEComponent       `json:"pue"`
}

type EnergyComponent struct {
	Units string  `json:"units"`
	Value float64 `json:"value"`
}

type EmissionsComponent struct {
	Region string  `json:"region"`
	Units  string  `json:"units"`
	Value  float64 `json:"value"`
}

type EmbodiedComponent struct {
	Units string  `json:"units"`
	Value float64 `json:"value"`
}

type PUEComponent struct {
	Provider string  `json:"provider"`
	Value    float64 `json:"value"`
}

// prometheusResponse is the Prometheus instant query HTTP API response envelope.
type prometheusResponse struct {
	Status string         `json:"status"`
	Error  string         `json:"error,omitempty"`
	Data   prometheusData `json:"data"`
}

type prometheusData struct {
	ResultType string             `json:"resultType"`
	Result     []prometheusResult `json:"result"`
}

type prometheusResult struct {
	Metric map[string]string `json:"metric"`
	// Value is [timestamp, "float_string"] per the Prometheus HTTP API.
	Value []json.RawMessage `json:"value"`
}

// ParsePrometheusResponse extracts the scalar float value from a Prometheus
// instant query HTTP API JSON response.
func ParsePrometheusResponse(response string) (float64, error) {
	var r prometheusResponse
	if err := json.Unmarshal([]byte(response), &r); err != nil {
		return 0, fmt.Errorf("parsing prometheus response: %w", err)
	}
	if r.Status != "success" {
		return 0, fmt.Errorf("prometheus query error: %s", r.Error)
	}
	if len(r.Data.Result) == 0 {
		return 0, fmt.Errorf("prometheus query returned no results")
	}
	result := r.Data.Result[0]
	if len(result.Value) != 2 {
		return 0, fmt.Errorf("unexpected prometheus value format: %v", result.Value)
	}
	// Value[0] is the timestamp (number), Value[1] is the metric value (quoted float).
	var valueStr string
	if err := json.Unmarshal(result.Value[1], &valueStr); err != nil {
		return 0, fmt.Errorf("parsing prometheus value: %w", err)
	}
	return strconv.ParseFloat(valueStr, 64)
}

// ComputeSCI computes the Software Carbon Intensity score from energy in Joules.
// Formula: SCI = (E * I * PUE) + M, where E is converted from Joules to kWh.
func ComputeSCI(energyJoules float64, cfg SCIConfig) SCI {
	energyKWh := energyJoules * wattSecondsToKWh
	sciValue := (energyKWh * cfg.EmissionsFactor.Value * cfg.PUE.Value) + cfg.EmbodiedCarbon.Value
	return SCI{
		Components: SCIComponents{
			EnergyConsumption: EnergyComponent{
				Units: "kWh",
				Value: energyKWh,
			},
			EmissionsFactor: cfg.EmissionsFactor,
			EmbodiedCarbon:  cfg.EmbodiedCarbon,
			PUE:             cfg.PUE,
		},
		Units: "gCO2eq",
		Value: sciValue,
	}
}
