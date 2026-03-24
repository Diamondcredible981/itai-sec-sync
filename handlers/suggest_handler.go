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
	Strategy         string      `json:"strategy"`
	Message          string      `json:"message,omitempty"`
	TotalOperations  int         `json:"total_operations"`
	Operations       []Operation `json:"operations"`
	AddOperations    []Operation `json:"add_operations"`
	RemoveOperations []Operation `json:"remove_operations"`
	Explanation      *SuggestExplanation `json:"explanation,omitempty"`
	Error            string      `json:"error,omitempty"`
}

func (s *Service) GetProductSuggestions(c *gin.Context) {
	strategy := c.DefaultQuery("strategy", "min-change")
	if strategy != "min-change" && strategy != "min-size" {
		s.badRequest(c, "strategy 参数仅支持 min-change 或 min-size")
		return
	}

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
			Strategy: strategy,
			Message: "无法覆盖所有功能，产品库不包含所有必需功能。",
			Error:   "无法覆盖",
		})
		return
	}

	// 核心逻辑重构：
	// 1) 以“修改次数最少”为第一目标
	// 2) 在同修改次数下，以“冗余最少”为第二目标
	// 3) 再以“组件总数更少”为第三目标
	bestTargetSet, minOperations, ok := findOptimalTargetSetByStrategy(strategy, allProducts, allFunctions, allProductIDs, allProductsMap, currentProductIDs, currentProductIDsSet, allFunctionIDsSet)
	if !ok {
		c.JSON(http.StatusOK, SuggestResponse{
			Success: false,
			Strategy: strategy,
			Message: "无法生成满足覆盖要求的优化方案。",
			Error:   "优化失败",
		})
		return
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
			Strategy:        strategy,
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
		Strategy:         strategy,
		Message:          "产品建议生成成功",
		TotalOperations:  minOperations,
		Operations:       append(addOps, removeOps...),
		AddOperations:    addOps,
		RemoveOperations: removeOps,
		Explanation:      explanation,
	})
}

func findOptimalTargetSetByStrategy(
	strategy string,
	allProducts []models.Product,
	allFunctions []models.Function,
	allProductIDs []int,
	allProductsMap map[int]models.Product,
	currentProductIDs []int,
	currentProductIDsSet map[int]bool,
	allFunctionIDsSet map[int]bool,
) ([]int, int, bool) {
	if strategy == "min-size" {
		return findOptimalTargetSetMinSize(allProductIDs, allProductsMap, allFunctions, allFunctionIDsSet, currentProductIDsSet)
	}
	return findOptimalTargetSet(allProducts, allFunctions, currentProductIDs)
}

func findOptimalTargetSetMinSize(
	allProductIDs []int,
	allProductsMap map[int]models.Product,
	allFunctions []models.Function,
	allFunctionIDsSet map[int]bool,
	currentProductIDsSet map[int]bool,
) ([]int, int, bool) {
	for k := 1; k <= len(allProductIDs); k++ {
		combinations := getCombinations(allProductIDs, k)
		bestFound := false
		bestOps := -1
		bestRedundant := -1
		bestAdd := -1
		var bestTarget []int

		for _, combo := range combinations {
			if !checkCoverage(combo, allFunctionIDsSet, allProductsMap) {
				continue
			}

			targetMetrics, _, _ := buildSuggestMetrics(combo, allFunctions, allProductsMap)
			totalOps, addOps := calcOps(combo, currentProductIDsSet)

			if !bestFound ||
				totalOps < bestOps ||
				(totalOps == bestOps && targetMetrics.RedundantCount < bestRedundant) ||
				(totalOps == bestOps && targetMetrics.RedundantCount == bestRedundant && addOps < bestAdd) {
				bestFound = true
				bestOps = totalOps
				bestRedundant = targetMetrics.RedundantCount
				bestAdd = addOps
				bestTarget = combo
			}
		}

		if bestFound {
			return bestTarget, bestOps, true
		}
	}

	return nil, -1, false
}

func calcOps(target []int, currentProductIDsSet map[int]bool) (int, int) {
	targetSet := make(map[int]bool, len(target))
	for _, id := range target {
		targetSet[id] = true
	}

	addOps := 0
	removeOps := 0
	for id := range targetSet {
		if !currentProductIDsSet[id] {
			addOps++
		}
	}
	for id := range currentProductIDsSet {
		if !targetSet[id] {
			removeOps++
		}
	}

	return addOps + removeOps, addOps
}

type suggestSearchState struct {
	Found       bool
	Redundant   int
	Size        int
	AddCount    int
	SelectedIDs []int
}

func findOptimalTargetSet(allProducts []models.Product, allFunctions []models.Function, currentProductIDs []int) ([]int, int, bool) {
	funcCount := len(allFunctions)
	if funcCount == 0 {
		return []int{}, 0, true
	}

	funcIndex := make(map[uint]int, funcCount)
	for i, f := range allFunctions {
		funcIndex[f.ID] = i
	}

	productIDs := make([]int, 0, len(allProducts))
	productFuncs := make([][]int, 0, len(allProducts))
	for _, p := range allProducts {
		productIDs = append(productIDs, int(p.ID))
		fIdxs := make([]int, 0)
		for _, fID := range utils.StringToIntSlice(p.FunctionIDsStr) {
			if idx, ok := funcIndex[uint(fID)]; ok {
				fIdxs = append(fIdxs, idx)
			}
		}
		productFuncs = append(productFuncs, fIdxs)
	}

	n := len(productIDs)
	if n == 0 {
		return nil, -1, false
	}

	currentSet := make(map[int]bool, len(currentProductIDs))
	for _, id := range currentProductIDs {
		currentSet[id] = true
	}
	isCurrent := make([]bool, n)
	for i, id := range productIDs {
		isCurrent[i] = currentSet[id]
	}

	wordCount := (funcCount + 63) / 64
	suffixUnion := make([][]uint64, n+1)
	suffixUnion[n] = make([]uint64, wordCount)
	for i := n - 1; i >= 0; i-- {
		words := make([]uint64, wordCount)
		copy(words, suffixUnion[i+1])
		for _, fIdx := range productFuncs[i] {
			w := fIdx / 64
			b := uint(fIdx % 64)
			words[w] |= uint64(1) << b
		}
		suffixUnion[i] = words
	}

	allMask := make([]uint64, wordCount)
	for i := 0; i < funcCount; i++ {
		w := i / 64
		b := uint(i % 64)
		allMask[w] |= uint64(1) << b
	}

	counts := make([]int, funcCount)
	coveredBits := make([]uint64, wordCount)
	selected := make([]bool, n)
	coveredCount := 0
	redundantCount := 0
	selectedSize := 0
	addCount := 0
	removeCount := 0

	var runSearch = func(opsLimit int) suggestSearchState {
		best := suggestSearchState{Found: false}

		var dfs func(idx int)
		dfs = func(idx int) {
			if addCount+removeCount > opsLimit {
				return
			}

			if !canStillCoverAll(coveredBits, suffixUnion[idx], allMask) {
				return
			}

			if idx == n {
				if coveredCount != funcCount {
					return
				}

				if !best.Found ||
					redundantCount < best.Redundant ||
					(redundantCount == best.Redundant && selectedSize < best.Size) ||
					(redundantCount == best.Redundant && selectedSize == best.Size && addCount < best.AddCount) {
					best.Found = true
					best.Redundant = redundantCount
					best.Size = selectedSize
					best.AddCount = addCount
					ids := make([]int, 0, selectedSize)
					for i, keep := range selected {
						if keep {
							ids = append(ids, productIDs[i])
						}
					}
					best.SelectedIDs = ids
				}
				return
			}

			tryStates := []bool{isCurrent[idx], !isCurrent[idx]}
			for _, include := range tryStates {
				selected[idx] = include

				deltaAdd, deltaRemove := 0, 0
				if include && !isCurrent[idx] {
					deltaAdd = 1
				}
				if !include && isCurrent[idx] {
					deltaRemove = 1
				}
				addCount += deltaAdd
				removeCount += deltaRemove

				deltaCovered := 0
				deltaRedundant := 0
				if include {
					selectedSize++
					for _, fIdx := range productFuncs[idx] {
						prev := counts[fIdx]
						counts[fIdx] = prev + 1
						if prev == 0 {
							deltaCovered++
							w := fIdx / 64
							b := uint(fIdx % 64)
							coveredBits[w] |= uint64(1) << b
						}
						if prev == 1 {
							deltaRedundant++
						}
					}
					coveredCount += deltaCovered
					redundantCount += deltaRedundant
				}

				dfs(idx + 1)

				if include {
					selectedSize--
					for _, fIdx := range productFuncs[idx] {
						prev := counts[fIdx]
						counts[fIdx] = prev - 1
						if prev == 1 {
							w := fIdx / 64
							b := uint(fIdx % 64)
							coveredBits[w] &^= uint64(1) << b
						}
					}
					coveredCount -= deltaCovered
					redundantCount -= deltaRedundant
				}

				addCount -= deltaAdd
				removeCount -= deltaRemove
				selected[idx] = false
			}
		}

		dfs(0)
		return best
	}

	for ops := 0; ops <= n; ops++ {
		best := runSearch(ops)
		if best.Found {
			return best.SelectedIDs, ops, true
		}
	}

	return nil, -1, false
}

func canStillCoverAll(coveredBits []uint64, suffixBits []uint64, allMask []uint64) bool {
	for i := range allMask {
		if (coveredBits[i]|suffixBits[i])&allMask[i] != allMask[i] {
			return false
		}
	}
	return true
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
