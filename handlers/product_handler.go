package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"github.com/iMayday-Yee/XinchuangAnalyze/utils"
)

// 获取所有产品
func (s *Service) ListProducts(c *gin.Context) {
	var products []models.Product

	// 支持按分类筛选
	typeID := c.Query("type_id")
	if typeID != "" {
		s.DB.Where("type_id = ?", typeID).Find(&products)
	} else {
		s.DB.Find(&products)
	}

	// 转换功能点字符串为数组
	for i := range products {
		products[i].FunctionIDs = utils.StringToIntSlice(products[i].FunctionIDsStr)
	}

	c.JSON(http.StatusOK, products)
}

// 获取单个产品
func (s *Service) GetProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	var product models.Product
	if err := s.DB.First(&product, id).Error; err != nil {
		s.notFound(c, "产品不存在")
		return
	}

	// 转换功能点字符串为数组
	product.FunctionIDs = utils.StringToIntSlice(product.FunctionIDsStr)

	c.JSON(http.StatusOK, product)
}

// 添加产品
func (s *Service) AddProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	// 验证产品类型是否存在
	var productType models.ProductType
	if err := s.DB.First(&productType, product.TypeID).Error; err != nil {
		s.badRequest(c, "产品类型不存在")
		return
	}

	// 去重功能点ID
	product.FunctionIDs = utils.UniqueIntSlice(product.FunctionIDs)

	// 验证功能点是否存在
	if len(product.FunctionIDs) > 0 {
		var existingFunctions []models.Function
		s.DB.Find(&existingFunctions)

		var validFunctionIDs []int
		for _, function := range existingFunctions {
			validFunctionIDs = append(validFunctionIDs, int(function.ID))
		}

		if !utils.ValidateIntSlice(product.FunctionIDs, validFunctionIDs) {
			s.badRequest(c, "部分功能点不存在")
			return
		}
	}

	// 转换功能点数组为字符串存储
	product.FunctionIDsStr = utils.IntSliceToString(product.FunctionIDs)

	if err := s.DB.Create(&product).Error; err != nil {
		s.internalError(c, "创建失败")
		return
	}

	c.JSON(http.StatusCreated, product)
}

// 更新产品
func (s *Service) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	var product models.Product
	if err := s.DB.First(&product, id).Error; err != nil {
		s.notFound(c, "产品不存在")
		return
	}

	var updateData models.Product
	if err := c.ShouldBindJSON(&updateData); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	// 验证产品类型是否存在
	if updateData.TypeID != 0 {
		var productType models.ProductType
		if err := s.DB.First(&productType, updateData.TypeID).Error; err != nil {
			s.badRequest(c, "产品类型不存在")
			return
		}
		product.TypeID = updateData.TypeID
	}

	// 验证功能点是否存在
	if len(updateData.FunctionIDs) > 0 {
		// 去重功能点ID
		updateData.FunctionIDs = utils.UniqueIntSlice(updateData.FunctionIDs)

		var existingFunctions []models.Function
		s.DB.Find(&existingFunctions)

		var validFunctionIDs []int
		for _, function := range existingFunctions {
			validFunctionIDs = append(validFunctionIDs, int(function.ID))
		}

		if !utils.ValidateIntSlice(updateData.FunctionIDs, validFunctionIDs) {
			s.badRequest(c, "部分功能点不存在")
			return
		}

		product.FunctionIDs = updateData.FunctionIDs
		product.FunctionIDsStr = utils.IntSliceToString(updateData.FunctionIDs)
	}

	// 更新其他字段
	if updateData.Name != "" {
		product.Name = updateData.Name
	}
	if updateData.Brand != "" {
		product.Brand = updateData.Brand
	}

	if err := s.DB.Save(&product).Error; err != nil {
		s.internalError(c, "更新失败")
		return
	}

	c.JSON(http.StatusOK, product)
}

// 删除产品
func (s *Service) DeleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	// 检查是否有拓扑节点引用此产品
	var usageCount int64
	s.DB.Model(&models.TopoNode{}).Where("product_id = ?", id).Count(&usageCount)
	if usageCount > 0 {
		s.badRequest(c, "该产品正在被网络拓扑使用，无法删除")
		return
	}

	if err := s.DB.Delete(&models.Product{}, id).Error; err != nil {
		s.internalError(c, "删除失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 批量获取产品（根据ID列表）
func (s *Service) GetProductsByIDs(c *gin.Context) {
	var request struct {
		ProductIDs []int `json:"product_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	// 去重产品ID
	request.ProductIDs = utils.UniqueIntSlice(request.ProductIDs)

	var products []models.Product
	s.DB.Where("id IN ?", request.ProductIDs).Find(&products)

	// 转换功能点字符串为数组
	for i := range products {
		products[i].FunctionIDs = utils.StringToIntSlice(products[i].FunctionIDsStr)
	}

	c.JSON(http.StatusOK, products)
}
