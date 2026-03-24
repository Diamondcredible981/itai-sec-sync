package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"github.com/iMayday-Yee/XinchuangAnalyze/utils"
)

// 获取所有网络拓扑
func (s *Service) ListTopos(c *gin.Context) {
	var topologies []models.NetworkTopo
	s.DB.Find(&topologies)

	// 转换产品ID字符串为数组，并获取产品详情
	for i := range topologies {
		topologies[i].ProductIDs = utils.StringToIntSlice(topologies[i].ProductIDsStr)

		// 获取关联的产品信息
		if len(topologies[i].ProductIDs) > 0 {
			var products []models.Product
			s.DB.Where("id IN ?", topologies[i].ProductIDs).Find(&products)
			for j := range products {
				products[j].FunctionIDs = utils.StringToIntSlice(products[j].FunctionIDsStr)
			}
			topologies[i].Products = products
		}
	}

	c.JSON(http.StatusOK, topologies)
}

// 获取单个网络拓扑
func (s *Service) GetTopo(c *gin.Context) {
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

	// 转换产品ID字符串为数组，并获取产品详情
	topology.ProductIDs = utils.StringToIntSlice(topology.ProductIDsStr)

	if len(topology.ProductIDs) > 0 {
		var products []models.Product
		s.DB.Where("id IN ?", topology.ProductIDs).Find(&products)
		for i := range products {
			products[i].FunctionIDs = utils.StringToIntSlice(products[i].FunctionIDsStr)
		}
		topology.Products = products
	}

	c.JSON(http.StatusOK, topology)
}

// 添加网络拓扑
func (s *Service) AddTopo(c *gin.Context) {
	var topology models.NetworkTopo
	if err := c.ShouldBindJSON(&topology); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	// 去重产品ID
	topology.ProductIDs = utils.UniqueIntSlice(topology.ProductIDs)

	// 验证产品是否存在
	if len(topology.ProductIDs) > 0 {
		var existingProducts []models.Product
		s.DB.Find(&existingProducts)

		var validProductIDs []int
		for _, product := range existingProducts {
			validProductIDs = append(validProductIDs, int(product.ID))
		}

		if !utils.ValidateIntSlice(topology.ProductIDs, validProductIDs) {
			s.badRequest(c, "部分产品不存在")
			return
		}
	}

	// 转换产品ID数组为字符串存储
	topology.ProductIDsStr = utils.IntSliceToString(topology.ProductIDs)

	if err := s.DB.Create(&topology).Error; err != nil {
		s.internalError(c, "创建失败")
		return
	}

	c.JSON(http.StatusCreated, topology)
}

// 更新网络拓扑
func (s *Service) UpdateTopo(c *gin.Context) {
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

	var updateData models.NetworkTopo
	if err := c.ShouldBindJSON(&updateData); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	// 验证产品是否存在
	if len(updateData.ProductIDs) > 0 {
		// 去重产品ID
		updateData.ProductIDs = utils.UniqueIntSlice(updateData.ProductIDs)

		var existingProducts []models.Product
		s.DB.Find(&existingProducts)

		var validProductIDs []int
		for _, product := range existingProducts {
			validProductIDs = append(validProductIDs, int(product.ID))
		}

		if !utils.ValidateIntSlice(updateData.ProductIDs, validProductIDs) {
			s.badRequest(c, "部分产品不存在")
			return
		}

		topology.ProductIDs = updateData.ProductIDs
		topology.ProductIDsStr = utils.IntSliceToString(updateData.ProductIDs)
	}

	// 更新其他字段
	if updateData.Name != "" {
		topology.Name = updateData.Name
	}

	if err := s.DB.Save(&topology).Error; err != nil {
		s.internalError(c, "更新失败")
		return
	}

	c.JSON(http.StatusOK, topology)
}

// 删除网络拓扑
func (s *Service) DeleteTopo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	if err := s.DB.Delete(&models.NetworkTopo{}, id).Error; err != nil {
		s.internalError(c, "删除失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 复制网络拓扑
func (s *Service) CopyTopo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	var originalTopo models.NetworkTopo
	if err := s.DB.First(&originalTopo, id).Error; err != nil {
		s.notFound(c, "网络拓扑不存在")
		return
	}

	// 创建副本
	newTopo := models.NetworkTopo{
		Name:          originalTopo.Name + "_副本",
		ProductIDsStr: originalTopo.ProductIDsStr,
	}

	if err := s.DB.Create(&newTopo).Error; err != nil {
		s.internalError(c, "复制失败")
		return
	}

	// 返回完整信息
	newTopo.ProductIDs = utils.StringToIntSlice(newTopo.ProductIDsStr)

	c.JSON(http.StatusCreated, newTopo)
}

// 获取拓扑的可视化数据
func (s *Service) GetTopoVisualization(c *gin.Context) {
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

	// 获取产品和产品类型信息
	topology.ProductIDs = utils.StringToIntSlice(topology.ProductIDsStr)

	var nodes []map[string]interface{}
	var edges []map[string]interface{}

	if len(topology.ProductIDs) > 0 {
		var products []models.Product
		s.DB.Where("id IN ?", topology.ProductIDs).Find(&products)

		// 获取产品类型信息
		var categoryMap = make(map[uint]models.ProductType)
		var categoryIDs []uint
		for _, product := range products {
			categoryIDs = append(categoryIDs, product.TypeID)
		}

		// 去重分类ID
		var uniqueCategoryIDs []int
		for _, cid := range categoryIDs {
			uniqueCategoryIDs = append(uniqueCategoryIDs, int(cid))
		}
		uniqueCategoryIDs = utils.UniqueIntSlice(uniqueCategoryIDs)

		var productTypes []models.ProductType
		s.DB.Where("id IN ?", uniqueCategoryIDs).Find(&productTypes)
		for _, pt := range productTypes {
			categoryMap[pt.ID] = pt
		}

		// 构建节点数据
		for i, product := range products {
			productType := categoryMap[product.TypeID]
			node := map[string]interface{}{
				"id":       strconv.Itoa(int(product.ID)),
				"label":    product.Name,
				"type":     productType.Name,
				"category": productType.Name,
				"icon":     productType.Icon,
				"color":    productType.Color,
				"x":        100 + (i%5)*150, // 简单的布局算法
				"y":        100 + (i/5)*100,
				"product":  product,
			}
			nodes = append(nodes, node)
		}

		// 构建简单的连接关系（这里可以根据实际需求优化）
		for i := 0; i < len(nodes)-1; i++ {
			edge := map[string]interface{}{
				"source": nodes[i]["id"],
				"target": nodes[i+1]["id"],
				"type":   "network",
			}
			edges = append(edges, edge)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"nodes":    nodes,
		"edges":    edges,
		"topology": topology,
	})
}
