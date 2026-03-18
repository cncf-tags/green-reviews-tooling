package sci

import (
	"math"
	"testing"
)

// TestComputeSCI verifies the SCI formula (E * I * PUE) + M against known inputs.
func TestComputeSCI(t *testing.T) {
	tests := []struct {
		name          string
		energyJoules  float64
		cfg           SCIConfig
		wantEnergyKWh float64
		wantSCI       float64
	}{
		{
			name:          "zero energy returns embodied carbon only",
			energyJoules:  0,
			cfg:           DefaultSCIConfig(),
			wantEnergyKWh: 0,
			wantSCI:       3.92, // M only
		},
		{
			name:          "default config matches grafana dashboard formula",
			energyJoules:  7200, // 0.002 kWh exactly
			cfg:           DefaultSCIConfig(),
			wantEnergyKWh: 0.002,
			// (0.002 * 78.81 * 1.42) + 3.92
			wantSCI: (0.002 * 78.81 * 1.42) + 3.92,
		},
		{
			name:         "custom config overrides defaults",
			energyJoules: 3600000, // 1 kWh exactly
			cfg: SCIConfig{
				EmissionsFactor: EmissionsComponent{Region: "US", Units: "gCO2eq/kWh", Value: 386.0},
				PUE:             PUEComponent{Provider: "Oracle", Value: 1.10},
				EmbodiedCarbon:  EmbodiedComponent{Units: "gCO2eq", Value: 5.00},
			},
			wantEnergyKWh: 1.0,
			// (1.0 * 386.0 * 1.10) + 5.0
			wantSCI: (1.0 * 386.0 * 1.10) + 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sci := ComputeSCI(tt.energyJoules, tt.cfg)

			if sci.Units != "gCO2eq" {
				t.Errorf("Units = %q, want %q", sci.Units, "gCO2eq")
			}
			if !almostEqual(sci.Components.EnergyConsumption.Value, tt.wantEnergyKWh, 1e-9) {
				t.Errorf("EnergyConsumption.Value = %v, want %v", sci.Components.EnergyConsumption.Value, tt.wantEnergyKWh)
			}
			if !almostEqual(sci.Value, tt.wantSCI, 1e-9) {
				t.Errorf("SCI.Value = %v, want %v", sci.Value, tt.wantSCI)
			}
			// Config fields should be reflected unchanged in the output.
			if sci.Components.EmissionsFactor.Region != tt.cfg.EmissionsFactor.Region {
				t.Errorf("EmissionsFactor.Region = %q, want %q", sci.Components.EmissionsFactor.Region, tt.cfg.EmissionsFactor.Region)
			}
			if sci.Components.PUE.Provider != tt.cfg.PUE.Provider {
				t.Errorf("PUE.Provider = %q, want %q", sci.Components.PUE.Provider, tt.cfg.PUE.Provider)
			}
		})
	}
}

// TestParsePrometheusResponse verifies happy-path JSON extraction.
func TestParsePrometheusResponse(t *testing.T) {
	response := `{
		"status": "success",
		"data": {
			"resultType": "vector",
			"result": [
				{
					"metric": {"container_name": "falco"},
					"value": [1698397200.123, "7867.507499998423"]
				}
			]
		}
	}`

	got, err := ParsePrometheusResponse(response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := 7867.507499998423
	if !almostEqual(got, want, 1e-6) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestParsePrometheusResponse_Empty verifies an error is returned for empty results.
func TestParsePrometheusResponse_Empty(t *testing.T) {
	response := `{
		"status": "success",
		"data": {
			"resultType": "vector",
			"result": []
		}
	}`

	_, err := ParsePrometheusResponse(response)
	if err == nil {
		t.Fatal("expected error for empty result set, got nil")
	}
}

// TestParsePrometheusResponse_Error verifies an error is returned for a Prometheus error status.
func TestParsePrometheusResponse_Error(t *testing.T) {
	response := `{
		"status": "error",
		"errorType": "bad_data",
		"error": "invalid query"
	}`

	_, err := ParsePrometheusResponse(response)
	if err == nil {
		t.Fatal("expected error for prometheus error status, got nil")
	}
}

// TestParsePrometheusResponse_Malformed verifies an error is returned for invalid JSON.
func TestParsePrometheusResponse_Malformed(t *testing.T) {
	_, err := ParsePrometheusResponse("not json")
	if err == nil {
		t.Fatal("expected error for malformed JSON, got nil")
	}
}

// almostEqual compares two floats within a tolerance.
func almostEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
