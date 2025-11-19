# Differential Privacy Anonymizer Plugin

A plugin that implements differential privacy mechanisms using Google's [Differential Privacy library](https://github.com/google/differential-privacy). This plugin provides privacy-preserving data anonymization through noise addition and differentially private aggregations.

## Features

This plugin implements the following differential privacy mechanisms:

### Noise Addition

- **Laplace Noise**: Adds calibrated Laplace noise to numeric values
- **Gaussian Noise**: Adds calibrated Gaussian noise to numeric values (provides (ε,δ)-differential privacy)

### Differentially Private Aggregations

- **Count**: Computes a differentially private count of values
- **Sum**: Computes a differentially private sum with bounded contributions
- **Mean**: Computes a differentially private mean with bounded contributions
- **Variance**: Computes a differentially private variance with bounded contributions

## Differential Privacy Concepts

### Key Parameters

1. **Epsilon (ε)**: Privacy budget

   - Smaller values = stronger privacy guarantees but more noise
   - Typical values: 0.1 to 10
   - Lower is better for privacy

2. **Delta (δ)**: Failure probability (for Gaussian mechanism)

   - Probability that privacy guarantee fails
   - Typical values: 1e-5 to 1e-10
   - Should be much smaller than 1/n where n is the dataset size

3. **Sensitivity**: Maximum change one individual can cause

   - For count: 1
   - For sum/mean/variance: depends on the bounds

4. **Bounds**: Min/max values for aggregations

   - Required for sum, mean, and variance
   - Used to calculate sensitivity and clip values

5. **Max Partitions Contributed**: Maximum number of partitions one user can contribute to
   - Affects the privacy budget allocation
   - Typically set to 1 for simple scenarios

## Usage Examples

### 1. Add Laplace Noise

```go
params := map[string]string{
    "value":       "100.0",
    "epsilon":     "1.0",
    "sensitivity": "1.0",
}
result, _ := plugin.Execute(ctx, "AddLaplaceNoise", params)
// Result: noisy value close to 100.0
```

### 2. Add Gaussian Noise

```go
params := map[string]string{
    "value":       "100.0",
    "epsilon":     "1.0",
    "delta":       "0.00001",
    "sensitivity": "1.0",
}
result, _ := plugin.Execute(ctx, "AddGaussianNoise", params)
// Result: noisy value close to 100.0
```

### 3. Differentially Private Count

```go
values := []float64{1, 2, 3, 4, 5}
valuesJSON, _ := json.Marshal(values)

params := map[string]string{
    "values":                     string(valuesJSON),
    "epsilon":                    "1.0",
    "delta":                      "0.00001",
    "max_partitions_contributed": "1",
}
result, _ := plugin.Execute(ctx, "DPCount", params)
// Result: noisy count close to 5
```

### 4. Differentially Private Sum

```go
values := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
valuesJSON, _ := json.Marshal(values)

params := map[string]string{
    "values":                     string(valuesJSON),
    "epsilon":                    "1.0",
    "delta":                      "0.00001",
    "lower_bound":                "0.0",
    "upper_bound":                "100.0",
    "max_partitions_contributed": "1",
}
result, _ := plugin.Execute(ctx, "DPSum", params)
// Result: noisy sum close to 150.0
```

### 5. Differentially Private Mean

```go
values := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
valuesJSON, _ := json.Marshal(values)

params := map[string]string{
    "values":                     string(valuesJSON),
    "epsilon":                    "1.0",
    "delta":                      "0.00001",
    "lower_bound":                "0.0",
    "upper_bound":                "100.0",
    "max_partitions_contributed": "1",
}
result, _ := plugin.Execute(ctx, "DPMean", params)
// Result: noisy mean close to 30.0
```

### 6. Differentially Private Variance

```go
values := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
valuesJSON, _ := json.Marshal(values)

params := map[string]string{
    "values":                     string(valuesJSON),
    "epsilon":                    "1.0",
    "delta":                      "0.00001",
    "lower_bound":                "0.0",
    "upper_bound":                "100.0",
    "max_partitions_contributed": "1",
}
result, _ := plugin.Execute(ctx, "DPVariance", params)
// Result: noisy variance close to actual variance
```

## Building the Plugin

```bash
# Build the plugin
make build-dpanonymizer

# Run tests
cd plugins/dpanonymizer
go test ./...

# Run example
cd example
go run main.go
```

## Use Cases

1. **Statistical Reporting**: Release aggregate statistics (counts, sums, means) while protecting individual privacy
2. **Data Analytics**: Perform privacy-preserving analysis on sensitive datasets
3. **Machine Learning**: Add noise to model outputs or training data
4. **Healthcare Analytics**: Analyze patient data while maintaining HIPAA compliance
5. **Financial Reporting**: Report financial metrics without revealing individual transactions

## Privacy Considerations

### Choosing Epsilon

- **High Privacy (ε < 1)**: Strong privacy but significant noise
- **Moderate Privacy (ε = 1-3)**: Balanced privacy and utility
- **Low Privacy (ε > 3)**: Weaker privacy but better accuracy

### Privacy Budget Management

- Each query consumes privacy budget (epsilon)
- Total epsilon across all queries should be limited
- Consider using privacy accounting for multiple queries

### Sensitivity Calculation

- For count: sensitivity = max_partitions_contributed
- For sum: sensitivity = max_partitions_contributed × (upper_bound - lower_bound)
- For mean/variance: similar to sum, but divided by count

## References

- [Google Differential Privacy Library](https://github.com/google/differential-privacy)
- [Differential Privacy Primer](https://programming-dp.com/)
- [The Algorithmic Foundations of Differential Privacy](https://www.cis.upenn.edu/~aaroth/Papers/privacybook.pdf)

## License

This plugin uses Google's Differential Privacy library, which is licensed under Apache License 2.0.
