# Differential Privacy Anonymizer Plugin - Usage Guide

## Quick Start

### 1. Build the Plugin

```bash
make build-dpanonymizer
```

This will create the plugin binary at `host-server/bin/plugins/dpanonymizer/dpanonymizer_v1.0.0`.

### 2. Run the Example

```bash
cd plugins/dpanonymizer/example
go run main.go
```

## Understanding Differential Privacy

Differential privacy is a mathematical framework for quantifying privacy loss. It ensures that the output of a computation doesn't reveal too much about any individual in the dataset.

### The Privacy Guarantee

A mechanism M provides (ε,δ)-differential privacy if for all datasets D1 and D2 that differ in one individual, and all possible outputs S:

```
Pr[M(D1) ∈ S] ≤ exp(ε) × Pr[M(D2) ∈ S] + δ
```

- **ε (epsilon)**: Privacy loss parameter. Smaller = better privacy
- **δ (delta)**: Failure probability. Should be very small (e.g., 1e-5)

## Practical Guide

### Scenario 1: Releasing a Count

**Problem**: You want to release the number of users who visited a website today.

**Solution**: Use DPCount

```go
// Suppose we have 1000 visits
values := make([]float64, 1000)
for i := range values {
    values[i] = 1.0
}
valuesJSON, _ := json.Marshal(values)

params := map[string]string{
    "values":                     string(valuesJSON),
    "epsilon":                    "0.5",  // Strong privacy
    "delta":                      "1e-5",
    "max_partitions_contributed": "1",    // Each user counted once
}

result, _ := plugin.Execute(ctx, "DPCount", params)
// Result will be close to 1000, but with noise added
```

**Why it works**: The noise masks whether any single user was in the dataset.

### Scenario 2: Computing Average Age

**Problem**: You want to compute the average age of users (ages range from 0-120).

**Solution**: Use DPMean

```go
ages := []float64{25, 30, 35, 40, 45, 50, 55, 60}
agesJSON, _ := json.Marshal(ages)

params := map[string]string{
    "values":                     string(agesJSON),
    "epsilon":                    "1.0",
    "delta":                      "1e-5",
    "lower_bound":                "0.0",   // Minimum age
    "upper_bound":                "120.0", // Maximum age
    "max_partitions_contributed": "1",     // Each user contributes once
}

result, _ := plugin.Execute(ctx, "DPMean", params)
// Result will be close to actual mean (42.5) with noise
```

**Important**:

- Bounds should reflect realistic min/max values
- Tighter bounds = less noise but values outside bounds get clipped

### Scenario 3: Total Revenue with Privacy

**Problem**: Release total revenue from sales, where each sale is between $0-$1000.

**Solution**: Use DPSum

```go
sales := []float64{100, 250, 500, 150, 300, 450, 200, 350}
salesJSON, _ := json.Marshal(sales)

params := map[string]string{
    "values":                     string(salesJSON),
    "epsilon":                    "2.0",     // Moderate privacy
    "delta":                      "1e-5",
    "lower_bound":                "0.0",
    "upper_bound":                "1000.0",
    "max_partitions_contributed": "1",
}

result, _ := plugin.Execute(ctx, "DPSum", params)
// Result will be close to 2300 with noise
```

### Scenario 4: Adding Noise to a Single Value

**Problem**: You have a sensitive metric (e.g., number of errors) and want to release it with privacy.

**Solution**: Use AddLaplaceNoise or AddGaussianNoise

```go
// Laplace mechanism (pure ε-DP)
params := map[string]string{
    "value":       "42.0",  // Actual error count
    "epsilon":     "1.0",
    "sensitivity": "1.0",   // One person can change count by at most 1
}
result, _ := plugin.Execute(ctx, "AddLaplaceNoise", params)

// Gaussian mechanism ((ε,δ)-DP, often less noise)
params := map[string]string{
    "value":       "42.0",
    "epsilon":     "1.0",
    "delta":       "1e-5",
    "sensitivity": "1.0",
}
result, _ := plugin.Execute(ctx, "AddGaussianNoise", params)
```

## Parameter Selection Guide

### Epsilon (Privacy Budget)

| Epsilon | Privacy Level | Use Case                                |
| ------- | ------------- | --------------------------------------- |
| 0.1     | Very Strong   | Highly sensitive data (medical records) |
| 0.5-1.0 | Strong        | Personal data (age, income)             |
| 1.0-3.0 | Moderate      | General analytics                       |
| > 3.0   | Weak          | Low-sensitivity data                    |

### Delta

- Should be much smaller than 1/n (where n = dataset size)
- Common values: 1e-5, 1e-6, 1e-10
- For dataset of 10,000 records: use δ ≤ 1e-5

### Bounds Selection

**Too Wide**: More noise, less accurate results
**Too Narrow**: Values get clipped, biased results

**Strategy**:

1. Analyze your data distribution
2. Set bounds to cover 95-99% of values
3. Accept that outliers will be clipped

Example:

```go
// If ages are typically 18-80, but max is 120
// Option 1: Tight bounds (less noise, some clipping)
"lower_bound": "18.0",
"upper_bound": "80.0",

// Option 2: Wide bounds (more noise, no clipping)
"lower_bound": "0.0",
"upper_bound": "120.0",
```

### Max Partitions Contributed

- How many times can one user appear in the aggregation?
- **Most common**: 1 (each user counted once)
- **Multiple contributions**: If users can have multiple records

## Common Pitfalls

### ❌ Pitfall 1: Reusing Privacy Budget

```go
// BAD: Using same epsilon for multiple queries
for i := 0; i < 10; i++ {
    result, _ := plugin.Execute(ctx, "DPCount", map[string]string{
        "epsilon": "1.0",  // Total privacy loss = 10.0!
        // ...
    })
}
```

**Solution**: Divide epsilon across queries or use privacy accounting.

### ❌ Pitfall 2: Ignoring Sensitivity

```go
// BAD: Wrong sensitivity for the query
params := map[string]string{
    "value":       "100.0",
    "epsilon":     "1.0",
    "sensitivity": "1.0",  // But value could change by 100!
}
```

**Solution**: Calculate sensitivity based on how much one individual can affect the result.

### ❌ Pitfall 3: Unrealistic Bounds

```go
// BAD: Bounds that don't match data
params := map[string]string{
    "lower_bound": "0.0",
    "upper_bound": "1000000.0",  // Way too wide for ages!
}
```

**Solution**: Use realistic bounds based on your data domain.

## Advanced Topics

### Privacy Budget Accounting

If you need to run multiple queries, track total epsilon:

```go
totalEpsilon := 0.0
epsilonPerQuery := 0.1
maxQueries := 10

for i := 0; i < maxQueries; i++ {
    if totalEpsilon + epsilonPerQuery > 1.0 {
        break  // Budget exhausted
    }

    // Run query with epsilonPerQuery
    totalEpsilon += epsilonPerQuery
}
```

### Composition Theorems

- **Sequential Composition**: Total ε = sum of individual εs
- **Parallel Composition**: Total ε = max of individual εs (if queries on disjoint data)

### When to Use Laplace vs Gaussian

**Laplace (Pure ε-DP)**:

- Simpler privacy guarantee
- No δ parameter
- More noise for same ε

**Gaussian ((ε,δ)-DP)**:

- Slightly weaker guarantee (due to δ)
- Less noise for same ε
- Better for multiple queries (advanced composition)

## Testing Your Implementation

```bash
# Run unit tests
cd plugins/dpanonymizer
go test ./impl -v

# Run example with different parameters
cd example
PLUGIN_PATH=../../../host-server/bin/plugins/dpanonymizer/dpanonymizer_v1.0.0 go run main.go
```

## Further Reading

1. **Beginner**: [A Friendly Introduction to Differential Privacy](https://desfontain.es/privacy/friendly-intro-to-differential-privacy.html)
2. **Intermediate**: [Programming Differential Privacy](https://programming-dp.com/)
3. **Advanced**: [The Algorithmic Foundations of Differential Privacy](https://www.cis.upenn.edu/~aaroth/Papers/privacybook.pdf)
4. **Google's DP Library**: [GitHub Repository](https://github.com/google/differential-privacy)
