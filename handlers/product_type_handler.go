package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
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
		s.badRequest(c, "无效的ID")
		return
	}

	var productType models.ProductType
	if err := s.DB.First(&productType, id).Error; err != nil {
		s.notFound(c, "产品类型不存在")
		return
	}

	c.JSON(http.StatusOK, productType)
}

// 添加产品类型
func (s *Service) AddProductType(c *gin.Context) {
	var productType models.ProductType
	if err := c.ShouldBindJSON(&productType); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	if err := s.DB.Create(&productType).Error; err != nil {
		s.internalError(c, "创建失败")
		return
	}

	c.JSON(http.StatusCreated, productType)
}

// 更新产品类型
func (s *Service) UpdateProductType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	var productType models.ProductType
	if err := s.DB.First(&productType, id).Error; err != nil {
		s.notFound(c, "产品类型不存在")
		return
	}

	if err := c.ShouldBindJSON(&productType); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	if err := s.DB.Save(&productType).Error; err != nil {
		s.internalError(c, "更新失败")
		return
	}

	c.JSON(http.StatusOK, productType)
}

// 删除产品类型
func (s *Service) DeleteProductType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	// 检查是否有产品使用此类型
	var count int64
	s.DB.Model(&models.Product{}).Where("type_id = ?", id).Count(&count)
	if count > 0 {
		s.badRequest(c, "该产品类型下还有产品，无法删除")
		return
	}

	if err := s.DB.Delete(&models.ProductType{}, id).Error; err != nil {
		s.internalError(c, "删除失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
