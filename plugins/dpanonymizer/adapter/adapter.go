package adapter

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/wylu1037/polyglot-plugin-showcase/plugins/dpanonymizer/impl"
	"github.com/wylu1037/polyglot-plugin-showcase/proto/common"
)

// DPAnonymizerAdapter adapts the differential privacy anonymizer implementation to the common plugin interface
type DPAnonymizerAdapter struct {
	impl *impl.DPAnonymizerImpl
}

func NewDPAnonymizerAdapter() *DPAnonymizerAdapter {
	return &DPAnonymizerAdapter{
		impl: &impl.DPAnonymizerImpl{},
	}
}

func (a *DPAnonymizerAdapter) GetMetadata() (*common.MetadataResponse, error) {
	return &common.MetadataResponse{
		Name:        "dpanonymizer",
		Version:     "1.0.0",
		Description: "Differential Privacy Anonymization plugin using Google's DP library",
		Methods: []string{
			"AddLaplaceNoise",
			"AddGaussianNoise",
			"DPCount",
			"DPSum",
			"DPMean",
			"DPVariance",
		},
		Capabilities: map[string]string{
			"type":              "anonymization",
			"privacy_mechanism": "differential_privacy",
			"noise_types":       "laplace,gaussian",
			"aggregations":      "count,sum,mean,variance",
		},
		ProtocolVersion: common.CurrentProtocolVersion,
	}, nil
}

func (a *DPAnonymizerAdapter) Execute(method string, params map[string]string) (*common.ExecuteResponse, error) {
	var result string
	var err error

	// Route to appropriate method
	switch method {
	case "AddLaplaceNoise":
		result, err = a.executeLaplaceNoise(params)
	case "AddGaussianNoise":
		result, err = a.executeGaussianNoise(params)
	case "DPCount":
		result, err = a.executeDPCount(params)
	case "DPSum":
		result, err = a.executeDPSum(params)
	case "DPMean":
		result, err = a.executeDPMean(params)
	case "DPVariance":
		result, err = a.executeDPVariance(params)
	default:
		errMsg := fmt.Sprintf("unknown method: %s", method)
		return &common.ExecuteResponse{
			Success: false,
			Error:   &errMsg,
		}, nil
	}

	if err != nil {
		errMsg := err.Error()
		return &common.ExecuteResponse{
			Success: false,
			Error:   &errMsg,
		}, nil
	}

	return &common.ExecuteResponse{
		Result:  &result,
		Success: true,
	}, nil
}

func (a *DPAnonymizerAdapter) executeLaplaceNoise(params map[string]string) (string, error) {
	value, err := strconv.ParseFloat(params["value"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid value parameter: %w", err)
	}

	epsilon, err := strconv.ParseFloat(params["epsilon"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid epsilon parameter: %w", err)
	}

	sensitivity, err := strconv.ParseFloat(params["sensitivity"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid sensitivity parameter: %w", err)
	}

	noisyValue, err := a.impl.AddLaplaceNoise(value, epsilon, sensitivity)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.6f", noisyValue), nil
}

func (a *DPAnonymizerAdapter) executeGaussianNoise(params map[string]string) (string, error) {
	value, err := strconv.ParseFloat(params["value"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid value parameter: %w", err)
	}

	epsilon, err := strconv.ParseFloat(params["epsilon"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid epsilon parameter: %w", err)
	}

	delta, err := strconv.ParseFloat(params["delta"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid delta parameter: %w", err)
	}

	sensitivity, err := strconv.ParseFloat(params["sensitivity"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid sensitivity parameter: %w", err)
	}

	noisyValue, err := a.impl.AddGaussianNoise(value, epsilon, delta, sensitivity)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.6f", noisyValue), nil
}

func (a *DPAnonymizerAdapter) executeDPCount(params map[string]string) (string, error) {
	valuesJSON, ok := params["values"]
	if !ok {
		return "", fmt.Errorf("missing values parameter")
	}

	var values []float64
	if err := json.Unmarshal([]byte(valuesJSON), &values); err != nil {
		return "", fmt.Errorf("invalid values parameter: %w", err)
	}

	epsilon, err := strconv.ParseFloat(params["epsilon"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid epsilon parameter: %w", err)
	}

	delta, err := strconv.ParseFloat(params["delta"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid delta parameter: %w", err)
	}

	maxPartitionsContributed, err := strconv.ParseInt(params["max_partitions_contributed"], 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid max_partitions_contributed parameter: %w", err)
	}

	count, err := a.impl.DPCount(values, epsilon, delta, maxPartitionsContributed)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", count), nil
}

func (a *DPAnonymizerAdapter) executeDPSum(params map[string]string) (string, error) {
	valuesJSON, ok := params["values"]
	if !ok {
		return "", fmt.Errorf("missing values parameter")
	}

	var values []float64
	if err := json.Unmarshal([]byte(valuesJSON), &values); err != nil {
		return "", fmt.Errorf("invalid values parameter: %w", err)
	}

	epsilon, err := strconv.ParseFloat(params["epsilon"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid epsilon parameter: %w", err)
	}

	delta, err := strconv.ParseFloat(params["delta"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid delta parameter: %w", err)
	}

	lowerBound, err := strconv.ParseFloat(params["lower_bound"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid lower_bound parameter: %w", err)
	}

	upperBound, err := strconv.ParseFloat(params["upper_bound"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid upper_bound parameter: %w", err)
	}

	maxPartitionsContributed, err := strconv.ParseInt(params["max_partitions_contributed"], 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid max_partitions_contributed parameter: %w", err)
	}

	sum, err := a.impl.DPSum(values, epsilon, delta, lowerBound, upperBound, maxPartitionsContributed)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.6f", sum), nil
}

func (a *DPAnonymizerAdapter) executeDPMean(params map[string]string) (string, error) {
	valuesJSON, ok := params["values"]
	if !ok {
		return "", fmt.Errorf("missing values parameter")
	}

	var values []float64
	if err := json.Unmarshal([]byte(valuesJSON), &values); err != nil {
		return "", fmt.Errorf("invalid values parameter: %w", err)
	}

	epsilon, err := strconv.ParseFloat(params["epsilon"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid epsilon parameter: %w", err)
	}

	delta, err := strconv.ParseFloat(params["delta"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid delta parameter: %w", err)
	}

	lowerBound, err := strconv.ParseFloat(params["lower_bound"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid lower_bound parameter: %w", err)
	}

	upperBound, err := strconv.ParseFloat(params["upper_bound"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid upper_bound parameter: %w", err)
	}

	maxPartitionsContributed, err := strconv.ParseInt(params["max_partitions_contributed"], 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid max_partitions_contributed parameter: %w", err)
	}

	mean, err := a.impl.DPMean(values, epsilon, delta, lowerBound, upperBound, maxPartitionsContributed)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.6f", mean), nil
}

func (a *DPAnonymizerAdapter) executeDPVariance(params map[string]string) (string, error) {
	valuesJSON, ok := params["values"]
	if !ok {
		return "", fmt.Errorf("missing values parameter")
	}

	var values []float64
	if err := json.Unmarshal([]byte(valuesJSON), &values); err != nil {
		return "", fmt.Errorf("invalid values parameter: %w", err)
	}

	epsilon, err := strconv.ParseFloat(params["epsilon"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid epsilon parameter: %w", err)
	}

	delta, err := strconv.ParseFloat(params["delta"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid delta parameter: %w", err)
	}

	lowerBound, err := strconv.ParseFloat(params["lower_bound"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid lower_bound parameter: %w", err)
	}

	upperBound, err := strconv.ParseFloat(params["upper_bound"], 64)
	if err != nil {
		return "", fmt.Errorf("invalid upper_bound parameter: %w", err)
	}

	maxPartitionsContributed, err := strconv.ParseInt(params["max_partitions_contributed"], 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid max_partitions_contributed parameter: %w", err)
	}

	variance, err := a.impl.DPVariance(values, epsilon, delta, lowerBound, upperBound, maxPartitionsContributed)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.6f", variance), nil
}
