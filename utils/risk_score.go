package utils

import (
	"math"
	"sort"
)

type RiskWeights struct {
	CoverageGap  float64 `json:"coverage_gap"`
	Redundancy   float64 `json:"redundancy"`
	SinglePoint  float64 `json:"single_point"`
}

type RiskScore struct {
	Score            float64     `json:"score"`
	Level            string      `json:"level"`
	CoverageGapRate  float64     `json:"coverage_gap_rate"`
	RedundancyRate   float64     `json:"redundancy_rate"`
	SinglePointRate  float64     `json:"single_point_rate"`
	TopContributors  []RiskContribution `json:"top_contributors"`
	Weights          RiskWeights `json:"weights"`
}

type RiskContribution struct {
	Dimension string  `json:"dimension"`
	Rate      float64 `json:"rate"`
	Weight    float64 `json:"weight"`
	Impact    float64 `json:"impact"`
	Reason    string  `json:"reason"`
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
			Level:           riskLevel(100),
			CoverageGapRate: 100,
			RedundancyRate:  0,
			SinglePointRate: 0,
			TopContributors: []RiskContribution{},
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

	contributions := []RiskContribution{
		{
			Dimension: "coverage_gap",
			Rate:      round2(coverageGapRate),
			Weight:    round2(weights.CoverageGap),
			Impact:    round2(coverageGapRate * weights.CoverageGap),
			Reason:    "未覆盖功能越多，暴露面越大",
		},
		{
			Dimension: "redundancy",
			Rate:      round2(redundancyRate),
			Weight:    round2(weights.Redundancy),
			Impact:    round2(redundancyRate * weights.Redundancy),
			Reason:    "冗余功能越多，配置复杂度越高",
		},
		{
			Dimension: "single_point",
			Rate:      round2(singlePointRate),
			Weight:    round2(weights.SinglePoint),
			Impact:    round2(singlePointRate * weights.SinglePoint),
			Reason:    "单点承载比例越高，抗失效能力越弱",
		},
	}

	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].Impact > contributions[j].Impact
	})

	return RiskScore{
		Score:           round2(score),
		Level:           riskLevel(score),
		CoverageGapRate: round2(coverageGapRate),
		RedundancyRate:  round2(redundancyRate),
		SinglePointRate: round2(singlePointRate),
		TopContributors: contributions,
		Weights:         weights,
	}
}

func riskLevel(score float64) string {
	switch {
	case score >= 75:
		return "critical"
	case score >= 50:
		return "high"
	case score >= 25:
		return "medium"
	default:
		return "low"
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
