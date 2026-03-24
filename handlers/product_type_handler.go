package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"net/http"
	"strconv"
)

// 获取所有产品类型
func (s *Service) ListProductTypes(c *gin.Context) {
	var productTypes []models.ProductType
	s.DB.Find(&productTypes)
	c.JSON(http.StatusOK, productTypes)
}

// 获取单个产品类型
func (s *Service) GetProductType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var productType models.ProductType
	if err := s.DB.First(&productType, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品类型不存在"})
		return
	}

	c.JSON(http.StatusOK, productType)
}

// 添加产品类型
func (s *Service) AddProductType(c *gin.Context) {
	var productType models.ProductType
	if err := c.ShouldBindJSON(&productType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.DB.Create(&productType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, productType)
}

// 更新产品类型
func (s *Service) UpdateProductType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var productType models.ProductType
	if err := s.DB.First(&productType, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "产品类型不存在"})
		return
	}

	if err := c.ShouldBindJSON(&productType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.DB.Save(&productType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, productType)
}

// 删除产品类型
func (s *Service) DeleteProductType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 检查是否有产品使用此类型
	var count int64
	s.DB.Model(&models.Product{}).Where("type_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该产品类型下还有产品，无法删除"})
		return
	}

	if err := s.DB.Delete(&models.ProductType{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
