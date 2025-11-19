// Package impl provides the implementation of the differential privacy anonymizer plugin.
package impl

import (
	"errors"
	"fmt"

	"github.com/google/differential-privacy/go/v3/dpagg"
	"github.com/google/differential-privacy/go/v3/noise"
)

// DPAnonymizerImpl is the concrete implementation of the DPAnonymizer interface.
type DPAnonymizerImpl struct{}

// AddLaplaceNoise adds Laplace noise to a numeric value for differential privacy.
// Parameters:
//   - value: the original numeric value
//   - epsilon: privacy budget (smaller = more privacy, more noise)
//   - sensitivity: the sensitivity of the query (how much one individual can affect the result)
func (d *DPAnonymizerImpl) AddLaplaceNoise(value, epsilon, sensitivity float64) (float64, error) {
	if epsilon <= 0 {
		return 0, errors.New("epsilon must be positive")
	}
	if sensitivity <= 0 {
		return 0, errors.New("sensitivity must be positive")
	}

	// Create Laplace noise generator
	ln := noise.Laplace()

	// Add noise to the value
	// l0sensitivity = 1 (one user contributes to one partition)
	// lInfSensitivity = sensitivity (how much one user can affect the result)
	// delta = 0 for Laplace (pure epsilon-DP)
	noisyValue, err := ln.AddNoiseFloat64(value, 1, sensitivity, epsilon, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to add Laplace noise: %w", err)
	}

	return noisyValue, nil
}

// AddGaussianNoise adds Gaussian noise to a numeric value for differential privacy.
// Parameters:
//   - value: the original numeric value
//   - epsilon: privacy budget
//   - delta: failure probability (typically very small, e.g., 1e-5)
//   - sensitivity: the sensitivity of the query
func (d *DPAnonymizerImpl) AddGaussianNoise(value, epsilon, delta, sensitivity float64) (float64, error) {
	if epsilon <= 0 {
		return 0, errors.New("epsilon must be positive")
	}
	if delta <= 0 || delta >= 1 {
		return 0, errors.New("delta must be between 0 and 1")
	}
	if sensitivity <= 0 {
		return 0, errors.New("sensitivity must be positive")
	}

	// Create Gaussian noise generator
	gn := noise.Gaussian()

	// Add noise to the value
	// l0sensitivity = 1 (one user contributes to one partition)
	// lInfSensitivity = sensitivity (how much one user can affect the result)
	noisyValue, err := gn.AddNoiseFloat64(value, 1, sensitivity, epsilon, delta)
	if err != nil {
		return 0, fmt.Errorf("failed to add Gaussian noise: %w", err)
	}

	return noisyValue, nil
}

// DPCount calculates a differentially private count of values.
func (d *DPAnonymizerImpl) DPCount(values []float64, epsilon, delta float64, maxPartitionsContributed int64) (int64, error) {
	if epsilon <= 0 {
		return 0, errors.New("epsilon must be positive")
	}
	if delta < 0 || delta >= 1 {
		return 0, errors.New("delta must be between 0 and 1")
	}
	if maxPartitionsContributed <= 0 {
		return 0, errors.New("maxPartitionsContributed must be positive")
	}

	// Create Count aggregator
	count, err := dpagg.NewCount(&dpagg.CountOptions{
		Epsilon:                  epsilon,
		Delta:                    delta,
		MaxPartitionsContributed: maxPartitionsContributed,
		Noise:                    noise.Gaussian(),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create Count aggregator: %w", err)
	}

	// Add each value (each represents one contribution)
	for range values {
		count.Increment()
	}

	// Get the noisy count result
	result, err := count.Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get count result: %w", err)
	}

	return result, nil
}

// DPSum calculates a differentially private sum of values.
func (d *DPAnonymizerImpl) DPSum(values []float64, epsilon, delta, lowerBound, upperBound float64, maxPartitionsContributed int64) (float64, error) {
	if epsilon <= 0 {
		return 0, errors.New("epsilon must be positive")
	}
	if delta < 0 || delta >= 1 {
		return 0, errors.New("delta must be between 0 and 1")
	}
	if maxPartitionsContributed <= 0 {
		return 0, errors.New("maxPartitionsContributed must be positive")
	}
	if lowerBound >= upperBound {
		return 0, errors.New("lowerBound must be less than upperBound")
	}

	// Create BoundedSumFloat64 aggregator
	sum, err := dpagg.NewBoundedSumFloat64(&dpagg.BoundedSumFloat64Options{
		Epsilon:                  epsilon,
		Delta:                    delta,
		MaxPartitionsContributed: maxPartitionsContributed,
		Lower:                    lowerBound,
		Upper:                    upperBound,
		Noise:                    noise.Gaussian(),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create BoundedSum aggregator: %w", err)
	}

	// Add all values
	for _, value := range values {
		if err := sum.Add(value); err != nil {
			return 0, fmt.Errorf("failed to add value: %w", err)
		}
	}

	// Get the noisy sum result
	result, err := sum.Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get sum result: %w", err)
	}

	return result, nil
}

// DPMean calculates a differentially private mean of values.
func (d *DPAnonymizerImpl) DPMean(values []float64, epsilon, delta, lowerBound, upperBound float64, maxPartitionsContributed int64) (float64, error) {
	if epsilon <= 0 {
		return 0, errors.New("epsilon must be positive")
	}
	if delta < 0 || delta >= 1 {
		return 0, errors.New("delta must be between 0 and 1")
	}
	if maxPartitionsContributed <= 0 {
		return 0, errors.New("maxPartitionsContributed must be positive")
	}
	if lowerBound >= upperBound {
		return 0, errors.New("lowerBound must be less than upperBound")
	}
	if len(values) == 0 {
		return 0, errors.New("values cannot be empty")
	}

	// Create BoundedMean aggregator
	mean, err := dpagg.NewBoundedMean(&dpagg.BoundedMeanOptions{
		Epsilon:                      epsilon,
		Delta:                        delta,
		MaxPartitionsContributed:     maxPartitionsContributed,
		MaxContributionsPerPartition: 1, // Each user contributes once per partition
		Lower:                        lowerBound,
		Upper:                        upperBound,
		Noise:                        noise.Gaussian(),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create BoundedMean aggregator: %w", err)
	}

	// Add all values
	for _, value := range values {
		if err := mean.Add(value); err != nil {
			return 0, fmt.Errorf("failed to add value: %w", err)
		}
	}

	// Get the noisy mean result
	result, err := mean.Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get mean result: %w", err)
	}

	return result, nil
}

// DPVariance calculates a differentially private variance of values.
func (d *DPAnonymizerImpl) DPVariance(values []float64, epsilon, delta, lowerBound, upperBound float64, maxPartitionsContributed int64) (float64, error) {
	if epsilon <= 0 {
		return 0, errors.New("epsilon must be positive")
	}
	if delta < 0 || delta >= 1 {
		return 0, errors.New("delta must be between 0 and 1")
	}
	if maxPartitionsContributed <= 0 {
		return 0, errors.New("maxPartitionsContributed must be positive")
	}
	if lowerBound >= upperBound {
		return 0, errors.New("lowerBound must be less than upperBound")
	}
	if len(values) == 0 {
		return 0, errors.New("values cannot be empty")
	}

	// Create BoundedVariance aggregator
	variance, err := dpagg.NewBoundedVariance(&dpagg.BoundedVarianceOptions{
		Epsilon:                      epsilon,
		Delta:                        delta,
		MaxPartitionsContributed:     maxPartitionsContributed,
		MaxContributionsPerPartition: 1, // Each user contributes once per partition
		Lower:                        lowerBound,
		Upper:                        upperBound,
		Noise:                        noise.Gaussian(),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create BoundedVariance aggregator: %w", err)
	}

	// Add all values
	for _, value := range values {
		if err := variance.Add(value); err != nil {
			return 0, fmt.Errorf("failed to add value: %w", err)
		}
	}

	// Get the noisy variance result
	result, err := variance.Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get variance result: %w", err)
	}

	return result, nil
}
