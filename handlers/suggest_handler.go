package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"github.com/iMayday-Yee/XinchuangAnalyze/utils"
)

type Operation struct {
	Type      string         `json:"type"` // "add" or "remove"
	ProductID int            `json:"product_id"`
	Product   models.Product `json:"product"`
	Reason    string         `json:"reason"`
}

type SuggestMetrics struct {
	CoverageRate   float64 `json:"coverage_rate"`
	RedundancyRate float64 `json:"redundancy_rate"`
	MissingCount   int     `json:"missing_count"`
	RedundantCount int     `json:"redundant_count"`
	RiskScore      float64 `json:"risk_score"`
}

type SuggestImprovement struct {
	CoverageRateDelta   float64 `json:"coverage_rate_delta"`
	RedundancyRateDelta float64 `json:"redundancy_rate_delta"`
	RiskScoreReduction  float64 `json:"risk_score_reduction"`
}

type SuggestExplanation struct {
	Current     SuggestMetrics     `json:"current"`
	Target      SuggestMetrics     `json:"target"`
	Improvement SuggestImprovement `json:"improvement"`
}

type SuggestResponse struct {
	Success          bool        `json:"success"`
	Message          string      `json:"message,omitempty"`
	TotalOperations  int         `json:"total_operations"`
	Operations       []Operation `json:"operations"`
	AddOperations    []Operation `json:"add_operations"`
	RemoveOperations []Operation `json:"remove_operations"`
	Explanation      *SuggestExplanation `json:"explanation,omitempty"`
	Error            string      `json:"error,omitempty"`
}

func (s *Service) GetProductSuggestions(c *gin.Context) {
	// 1. 获取拓扑ID
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	// 2. 获取拓扑信息
	var topology models.NetworkTopo
	if err := s.DB.First(&topology, id).Error; err != nil {
		s.notFound(c, "网络拓扑不存在")
		return
	}
	topology.ProductIDs = utils.StringToIntSlice(topology.ProductIDsStr)

	// 3. 获取所有功能点和所有产品
	var allFunctions []models.Function
	s.DB.Find(&allFunctions)
	var allProducts []models.Product
	s.DB.Find(&allProducts)

	// 4. 数据准备
	allProductsMap := make(map[int]models.Product)
	var allProductIDs []int
	for _, p := range allProducts {
		product := p // 必须在循环内创建新变量
		product.FunctionIDs = utils.StringToIntSlice(product.FunctionIDsStr)
		allProductsMap[int(product.ID)] = product
		allProductIDs = append(allProductIDs, int(product.ID))
	}

	allFunctionIDsSet := make(map[int]bool)
	functionNameMap := make(map[int]string)
	for _, f := range allFunctions {
		allFunctionIDsSet[int(f.ID)] = true
		functionNameMap[int(f.ID)] = f.Name
	}

	currentProductIDs := topology.ProductIDs
	currentProductIDsSet := make(map[int]bool)
	for _, pID := range currentProductIDs {
		currentProductIDsSet[pID] = true
	}

	currentMetrics, _, missingCurrentSet := buildSuggestMetrics(currentProductIDs, allFunctions, allProductsMap)

	// 5. 检查所有产品集合是否能覆盖所有功能点
	if !checkCoverage(allProductIDs, allFunctionIDsSet, allProductsMap) {
		c.JSON(http.StatusOK, SuggestResponse{
			Success: false,
			Message: "无法覆盖所有功能，产品库不包含所有必需功能。",
			Error:   "无法覆盖",
		})
		return
	}

	// 核心逻辑重构：
	// 1. 找到能覆盖所有功能的最小设备集尺寸k
	// 2. 在所有尺寸为k的设备集中，找到与当前设备集操作距离最近的一个
	// 3. 生成该操作的建议

	// 步骤1: 找到最小设备集尺寸k
	minProductSetSize := -1
	for k := 1; k <= len(allProductIDs); k++ {
		combinations := getCombinations(allProductIDs, k)
		for _, combo := range combinations {
			if checkCoverage(combo, allFunctionIDsSet, allProductsMap) {
				minProductSetSize = k
				break
			}
		}
		if minProductSetSize != -1 {
			break
		}
	}

	// 步骤2: 寻找最优目标集
	var bestTargetSet []int
	minOperations := -1

	optimalSizeCombinations := getCombinations(allProductIDs, minProductSetSize)

	for _, targetSetSlice := range optimalSizeCombinations {
		if !checkCoverage(targetSetSlice, allFunctionIDsSet, allProductsMap) {
			continue
		}

		targetSetMap := make(map[int]bool)
		for _, pID := range targetSetSlice {
			targetSetMap[pID] = true
		}

		// 计算操作数
		adds, removes := 0, 0
		for pID := range targetSetMap {
			if !currentProductIDsSet[pID] {
				adds++
			}
		}
		for pID := range currentProductIDsSet {
			if !targetSetMap[pID] {
				removes++
			}
		}
		totalOps := adds + removes

		if minOperations == -1 || totalOps < minOperations {
			minOperations = totalOps
			bestTargetSet = targetSetSlice
		}
	}

	// 步骤3: 生成结果
	if minOperations == 0 {
		explanation := &SuggestExplanation{
			Current: currentMetrics,
			Target:  currentMetrics,
			Improvement: SuggestImprovement{
				CoverageRateDelta:   0,
				RedundancyRateDelta: 0,
				RiskScoreReduction:  0,
			},
		}

		c.JSON(http.StatusOK, SuggestResponse{
			Success:         true,
			Message:         "当前配置已覆盖所有功能点，且无冗余，无需调整。",
			TotalOperations: 0,
			Explanation:     explanation,
		})
		return
	}

	bestTargetSetMap := make(map[int]bool)
	for _, pID := range bestTargetSet {
		bestTargetSetMap[pID] = true
	}

	addOps := make([]Operation, 0)
	for _, pID := range bestTargetSet {
		if !currentProductIDsSet[pID] {
			reason := buildAddReason(allProductsMap[pID], missingCurrentSet, functionNameMap)
			addOps = append(addOps, Operation{
				Type:      "add",
				ProductID: pID,
				Product:   allProductsMap[pID],
				Reason:    reason,
			})
		}
	}

	targetMetrics, targetFunctionMap, _ := buildSuggestMetrics(bestTargetSet, allFunctions, allProductsMap)

	removeOps := make([]Operation, 0)
	for _, pID := range currentProductIDs {
		if !bestTargetSetMap[pID] {
			reason := buildRemoveReason(allProductsMap[pID], targetFunctionMap, functionNameMap)
			removeOps = append(removeOps, Operation{
				Type:      "remove",
				ProductID: pID,
				Product:   allProductsMap[pID],
				Reason:    reason,
			})
		}
	}

	explanation := &SuggestExplanation{
		Current: currentMetrics,
		Target:  targetMetrics,
		Improvement: SuggestImprovement{
			CoverageRateDelta:   round2(targetMetrics.CoverageRate - currentMetrics.CoverageRate),
			RedundancyRateDelta: round2(targetMetrics.RedundancyRate - currentMetrics.RedundancyRate),
			RiskScoreReduction:  round2(currentMetrics.RiskScore - targetMetrics.RiskScore),
		},
	}

	c.JSON(http.StatusOK, SuggestResponse{
		Success:          true,
		Message:          "产品建议生成成功",
		TotalOperations:  minOperations,
		Operations:       append(addOps, removeOps...),
		AddOperations:    addOps,
		RemoveOperations: removeOps,
		Explanation:      explanation,
	})
}

func buildSuggestMetrics(productIDs []int, allFunctions []models.Function, allProductsMap map[int]models.Product) (SuggestMetrics, map[uint]int, map[int]bool) {
	functionMap := make(map[uint]int)
	for _, function := range allFunctions {
		functionMap[function.ID] = 0
	}

	for _, pID := range productIDs {
		if product, ok := allProductsMap[pID]; ok {
			for _, fID := range product.FunctionIDs {
				functionMap[uint(fID)]++
			}
		}
	}

	missingCount := 0
	redundantCount := 0
	missingSet := make(map[int]bool)

	for _, function := range allFunctions {
		count := functionMap[function.ID]
		if count == 0 {
			missingCount++
			missingSet[int(function.ID)] = true
		}
		if count > 1 {
			redundantCount++
		}
	}

	total := len(allFunctions)
	coverageRate := 0.0
	redundancyRate := 0.0
	if total > 0 {
		coverageRate = float64(total-missingCount) / float64(total) * 100
		redundancyRate = float64(redundantCount) / float64(total) * 100
	}

	risk := utils.CalculateRiskScore(total, functionMap, utils.DefaultRiskWeights())

	return SuggestMetrics{
		CoverageRate:   round2(coverageRate),
		RedundancyRate: round2(redundancyRate),
		MissingCount:   missingCount,
		RedundantCount: redundantCount,
		RiskScore:      risk.Score,
	}, functionMap, missingSet
}

func buildAddReason(product models.Product, missingCurrentSet map[int]bool, functionNameMap map[int]string) string {
	matched := make([]string, 0)
	seen := make(map[string]bool)

	for _, fID := range product.FunctionIDs {
		if missingCurrentSet[fID] {
			name := functionNameMap[fID]
			if name != "" && !seen[name] {
				matched = append(matched, name)
				seen[name] = true
			}
		}
	}

	if len(matched) == 0 {
		return "增强弹性并靠近最优配置"
	}

	return "补齐缺失功能: " + strings.Join(limitNames(matched, 3), "、")
}

func buildRemoveReason(product models.Product, targetFunctionMap map[uint]int, functionNameMap map[int]string) string {
	safeCovered := make([]string, 0)
	seen := make(map[string]bool)

	for _, fID := range product.FunctionIDs {
		if targetFunctionMap[uint(fID)] > 0 {
			name := functionNameMap[fID]
			if name != "" && !seen[name] {
				safeCovered = append(safeCovered, name)
				seen[name] = true
			}
		}
	}

	if len(safeCovered) == 0 {
		return "移除低价值设备，降低复杂度"
	}

	return "移除后仍可覆盖功能: " + strings.Join(limitNames(safeCovered, 3), "、")
}

func limitNames(names []string, limit int) []string {
	if len(names) <= limit {
		return names
	}
	return names[:limit]
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// checkCoverage 检查给定的产品ID列表是否能覆盖所有功能点
func checkCoverage(productIDs []int, allFunctionIDs map[int]bool, allProductsMap map[int]models.Product) bool {
	coveredFunctions := make(map[int]bool)
	for _, pID := range productIDs {
		if product, ok := allProductsMap[pID]; ok {
			for _, fID := range product.FunctionIDs {
				coveredFunctions[fID] = true
			}
		}
	}
	return len(coveredFunctions) >= len(allFunctionIDs)
}

// getCombinations 获取一个切片中k个元素的所有组合
func getCombinations(items []int, k int) [][]int {
	if k == 0 {
		return [][]int{{}}
	}
	if k < 0 || k > len(items) {
		return [][]int{}
	}

	var result [][]int
	var backtrack func(start int, combo []int)
	backtrack = func(start int, combo []int) {
		if len(combo) == k {
			newCombo := make([]int, k)
			copy(newCombo, combo)
			result = append(result, newCombo)
			return
		}

		for i := start; i < len(items); i++ {
			combo = append(combo, items[i])
			backtrack(i+1, combo)
			combo = combo[:len(combo)-1] // 回溯
		}
	}

	backtrack(0, []int{})
	return result
}
