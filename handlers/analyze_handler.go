package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"github.com/iMayday-Yee/XinchuangAnalyze/utils"
)

type AnalyzeResult struct {
	Redundant      []models.Function `json:"redundant"`
	Missing        []models.Function `json:"missing"`
	CoverageRate   float64           `json:"coverage_rate"`
	RedundancyRate float64           `json:"redundancy_rate"`
	Risk           utils.RiskScore   `json:"risk"`
}

func (s *Service) AnalyzeByProductIDs(c *gin.Context) {
	var topology models.NetworkTopo
	if err := c.ShouldBindJSON(&topology); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, s.buildAnalyzeResult(topology.ProductIDs))
}

func (s *Service) AnalyzeByTopoID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}
	var topology models.NetworkTopo
	if err := s.DB.First(&topology, id).Error; err != nil {
		s.notFound(c, "网络拓扑不存在")
		return
	}
	topology.ProductIDs = utils.StringToIntSlice(topology.ProductIDsStr)

	c.JSON(http.StatusOK, s.buildAnalyzeResult(topology.ProductIDs))
}

func (s *Service) buildAnalyzeResult(productIDs []int) AnalyzeResult {
	var allFunctions []models.Function
	s.DB.Find(&allFunctions)

	functionMap := make(map[uint]int)
	for _, function := range allFunctions {
		functionMap[function.ID] = 0
	}

	for _, productID := range productIDs {
		var product models.Product
		s.DB.First(&product, productID)
		functionIDs := utils.StringToIntSlice(product.FunctionIDsStr)
		for _, functionID := range functionIDs {
			functionMap[uint(functionID)]++
		}
	}

	var redundantFunctions []models.Function
	var missingFunctions []models.Function

	for _, function := range allFunctions {
		count := functionMap[function.ID]
		if count > 1 {
			redundantFunctions = append(redundantFunctions, function)
		}
		if count == 0 {
			missingFunctions = append(missingFunctions, function)
		}
	}

	total := len(allFunctions)
	coverageRate := 0.0
	redundancyRate := 0.0
	if total > 0 {
		coverageRate = float64(total-len(missingFunctions)) / float64(total) * 100
		redundancyRate = float64(len(redundantFunctions)) / float64(total) * 100
	}

	risk := utils.CalculateRiskScore(total, functionMap, utils.DefaultRiskWeights())

	return AnalyzeResult{
		Redundant:      redundantFunctions,
		Missing:        missingFunctions,
		CoverageRate:   coverageRate,
		RedundancyRate: redundancyRate,
		Risk:           risk,
	}
}
