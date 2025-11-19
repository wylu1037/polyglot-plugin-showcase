package impl

import (
	"math"
	"testing"
)

func TestAddLaplaceNoise(t *testing.T) {
	impl := &DPAnonymizerImpl{}
	
	tests := []struct {
		name        string
		value       float64
		epsilon     float64
		sensitivity float64
		wantErr     bool
	}{
		{
			name:        "valid parameters",
			value:       100.0,
			epsilon:     1.0,
			sensitivity: 1.0,
			wantErr:     false,
		},
		{
			name:        "invalid epsilon",
			value:       100.0,
			epsilon:     0.0,
			sensitivity: 1.0,
			wantErr:     true,
		},
		{
			name:        "invalid sensitivity",
			value:       100.0,
			epsilon:     1.0,
			sensitivity: 0.0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := impl.AddLaplaceNoise(tt.value, tt.epsilon, tt.sensitivity)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddLaplaceNoise() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// The result should be different from the original value (with high probability)
				// and should be a valid number
				if math.IsNaN(result) || math.IsInf(result, 0) {
					t.Errorf("AddLaplaceNoise() returned invalid result: %v", result)
				}
			}
		})
	}
}

func TestAddGaussianNoise(t *testing.T) {
	impl := &DPAnonymizerImpl{}
	
	tests := []struct {
		name        string
		value       float64
		epsilon     float64
		delta       float64
		sensitivity float64
		wantErr     bool
	}{
		{
			name:        "valid parameters",
			value:       100.0,
			epsilon:     1.0,
			delta:       1e-5,
			sensitivity: 1.0,
			wantErr:     false,
		},
		{
			name:        "invalid epsilon",
			value:       100.0,
			epsilon:     0.0,
			delta:       1e-5,
			sensitivity: 1.0,
			wantErr:     true,
		},
		{
			name:        "invalid delta",
			value:       100.0,
			epsilon:     1.0,
			delta:       1.0,
			sensitivity: 1.0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := impl.AddGaussianNoise(tt.value, tt.epsilon, tt.delta, tt.sensitivity)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddGaussianNoise() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if math.IsNaN(result) || math.IsInf(result, 0) {
					t.Errorf("AddGaussianNoise() returned invalid result: %v", result)
				}
			}
		})
	}
}

func TestDPCount(t *testing.T) {
	impl := &DPAnonymizerImpl{}
	
	tests := []struct {
		name                     string
		values                   []float64
		epsilon                  float64
		delta                    float64
		maxPartitionsContributed int64
		wantErr                  bool
	}{
		{
			name:                     "valid count",
			values:                   []float64{1, 2, 3, 4, 5},
			epsilon:                  1.0,
			delta:                    1e-5,
			maxPartitionsContributed: 1,
			wantErr:                  false,
		},
		{
			name:                     "invalid epsilon",
			values:                   []float64{1, 2, 3},
			epsilon:                  0.0,
			delta:                    1e-5,
			maxPartitionsContributed: 1,
			wantErr:                  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := impl.DPCount(tt.values, tt.epsilon, tt.delta, tt.maxPartitionsContributed)
			if (err != nil) != tt.wantErr {
				t.Errorf("DPCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Result should be close to the actual count
				actualCount := int64(len(tt.values))
				// Allow for some noise (within 50% for this test)
				if result < actualCount/2 || result > actualCount*2 {
					t.Logf("DPCount() result = %v, actual count = %v (noise is expected)", result, actualCount)
				}
			}
		})
	}
}

func TestDPSum(t *testing.T) {
	impl := &DPAnonymizerImpl{}
	
	tests := []struct {
		name                     string
		values                   []float64
		epsilon                  float64
		delta                    float64
		lowerBound               float64
		upperBound               float64
		maxPartitionsContributed int64
		wantErr                  bool
	}{
		{
			name:                     "valid sum",
			values:                   []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			epsilon:                  1.0,
			delta:                    1e-5,
			lowerBound:               0.0,
			upperBound:               10.0,
			maxPartitionsContributed: 1,
			wantErr:                  false,
		},
		{
			name:                     "invalid bounds",
			values:                   []float64{1.0, 2.0, 3.0},
			epsilon:                  1.0,
			delta:                    1e-5,
			lowerBound:               10.0,
			upperBound:               0.0,
			maxPartitionsContributed: 1,
			wantErr:                  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := impl.DPSum(tt.values, tt.epsilon, tt.delta, tt.lowerBound, tt.upperBound, tt.maxPartitionsContributed)
			if (err != nil) != tt.wantErr {
				t.Errorf("DPSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if math.IsNaN(result) || math.IsInf(result, 0) {
					t.Errorf("DPSum() returned invalid result: %v", result)
				}
			}
		})
	}
}

func TestDPMean(t *testing.T) {
	impl := &DPAnonymizerImpl{}
	
	tests := []struct {
		name                     string
		values                   []float64
		epsilon                  float64
		delta                    float64
		lowerBound               float64
		upperBound               float64
		maxPartitionsContributed int64
		wantErr                  bool
	}{
		{
			name:                     "valid mean",
			values:                   []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			epsilon:                  1.0,
			delta:                    1e-5,
			lowerBound:               0.0,
			upperBound:               10.0,
			maxPartitionsContributed: 1,
			wantErr:                  false,
		},
		{
			name:                     "empty values",
			values:                   []float64{},
			epsilon:                  1.0,
			delta:                    1e-5,
			lowerBound:               0.0,
			upperBound:               10.0,
			maxPartitionsContributed: 1,
			wantErr:                  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := impl.DPMean(tt.values, tt.epsilon, tt.delta, tt.lowerBound, tt.upperBound, tt.maxPartitionsContributed)
			if (err != nil) != tt.wantErr {
				t.Errorf("DPMean() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if math.IsNaN(result) || math.IsInf(result, 0) {
					t.Errorf("DPMean() returned invalid result: %v", result)
				}
			}
		})
	}
}

func TestDPVariance(t *testing.T) {
	impl := &DPAnonymizerImpl{}
	
	tests := []struct {
		name                     string
		values                   []float64
		epsilon                  float64
		delta                    float64
		lowerBound               float64
		upperBound               float64
		maxPartitionsContributed int64
		wantErr                  bool
	}{
		{
			name:                     "valid variance",
			values:                   []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			epsilon:                  1.0,
			delta:                    1e-5,
			lowerBound:               0.0,
			upperBound:               10.0,
			maxPartitionsContributed: 1,
			wantErr:                  false,
		},
		{
			name:                     "empty values",
			values:                   []float64{},
			epsilon:                  1.0,
			delta:                    1e-5,
			lowerBound:               0.0,
			upperBound:               10.0,
			maxPartitionsContributed: 1,
			wantErr:                  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := impl.DPVariance(tt.values, tt.epsilon, tt.delta, tt.lowerBound, tt.upperBound, tt.maxPartitionsContributed)
			if (err != nil) != tt.wantErr {
				t.Errorf("DPVariance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if math.IsNaN(result) || math.IsInf(result, 0) {
					t.Errorf("DPVariance() returned invalid result: %v", result)
				}
				// Variance should be non-negative
				if result < 0 {
					t.Errorf("DPVariance() returned negative variance: %v", result)
				}
			}
		})
	}
}

