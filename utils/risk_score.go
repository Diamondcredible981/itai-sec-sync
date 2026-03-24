package utils

import "math"

type RiskWeights struct {
	CoverageGap  float64 `json:"coverage_gap"`
	Redundancy   float64 `json:"redundancy"`
	SinglePoint  float64 `json:"single_point"`
}

type RiskScore struct {
	Score            float64     `json:"score"`
	CoverageGapRate  float64     `json:"coverage_gap_rate"`
	RedundancyRate   float64     `json:"redundancy_rate"`
	SinglePointRate  float64     `json:"single_point_rate"`
	Weights          RiskWeights `json:"weights"`
}

func DefaultRiskWeights() RiskWeights {
	return RiskWeights{
		CoverageGap: 0.5,
		Redundancy:  0.3,
		SinglePoint: 0.2,
	}
}

func CalculateRiskScore(totalFunctions int, functionMap map[uint]int, weights RiskWeights) RiskScore {
	if totalFunctions <= 0 {
		return RiskScore{
			Score:           100,
			CoverageGapRate: 100,
			RedundancyRate:  0,
			SinglePointRate: 0,
			Weights:         weights,
		}
	}

	missing := 0
	redundant := 0
	singlePoint := 0

	for _, count := range functionMap {
		if count == 0 {
			missing++
		}
		if count > 1 {
			redundant++
		}
		if count == 1 {
			singlePoint++
		}
	}

	coverageGapRate := float64(missing) / float64(totalFunctions) * 100
	redundancyRate := float64(redundant) / float64(totalFunctions) * 100
	singlePointRate := float64(singlePoint) / float64(totalFunctions) * 100

	score := coverageGapRate*weights.CoverageGap +
		redundancyRate*weights.Redundancy +
		singlePointRate*weights.SinglePoint

	score = clamp(score, 0, 100)

	return RiskScore{
		Score:           round2(score),
		CoverageGapRate: round2(coverageGapRate),
		RedundancyRate:  round2(redundancyRate),
		SinglePointRate: round2(singlePointRate),
		Weights:         weights,
	}
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
