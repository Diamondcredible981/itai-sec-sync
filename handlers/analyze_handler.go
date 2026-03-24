package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"github.com/iMayday-Yee/XinchuangAnalyze/utils"
	"net/http"
	"strconv"
)

func (s *Service) AnalyzeByProductIDs(c *gin.Context) {
	var topology models.NetworkTopo
	if err := c.ShouldBindJSON(&topology); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取所有功能点
	var allFunctions []models.Function
	s.DB.Find(&allFunctions)

	functionMap := make(map[uint]int)
	for _, function := range allFunctions {
		functionMap[function.ID] = 0
	}

	// 统计各功能点的覆盖情况
	for _, productID := range topology.ProductIDs {
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

	c.JSON(http.StatusOK, gin.H{
		"redundant":       redundantFunctions,
		"missing":         missingFunctions,
		"coverage_rate":   float64(len(allFunctions)-len(missingFunctions)) / float64(len(allFunctions)) * 100,
		"redundancy_rate": float64(len(redundantFunctions)) / float64(len(allFunctions)) * 100,
	})
}

func (s *Service) AnalyzeByTopoID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}
	var topology models.NetworkTopo
	if err := s.DB.First(&topology, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "网络拓扑不存在"})
		return
	}
	topology.ProductIDs = utils.StringToIntSlice(topology.ProductIDsStr)

	// 获取所有功能点
	var allFunctions []models.Function
	s.DB.Find(&allFunctions)

	functionMap := make(map[uint]int)
	for _, function := range allFunctions {
		functionMap[function.ID] = 0
	}

	// 统计各功能点的覆盖情况
	for _, productID := range topology.ProductIDs {
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

	c.JSON(http.StatusOK, gin.H{
		"redundant":       redundantFunctions,
		"missing":         missingFunctions,
		"coverage_rate":   float64(len(allFunctions)-len(missingFunctions)) / float64(len(allFunctions)) * 100,
		"redundancy_rate": float64(len(redundantFunctions)) / float64(len(allFunctions)) * 100,
	})
}
