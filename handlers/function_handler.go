package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
)

// 获取所有功能点
func (s *Service) ListFunctions(c *gin.Context) {
	var functions []models.Function
	s.DB.Find(&functions)
	c.JSON(http.StatusOK, functions)
}

// 获取单个功能点
func (s *Service) GetFunction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	var function models.Function
	if err := s.DB.First(&function, id).Error; err != nil {
		s.notFound(c, "功能点不存在")
		return
	}

	c.JSON(http.StatusOK, function)
}

// 添加功能点
func (s *Service) AddFunction(c *gin.Context) {
	var function models.Function
	if err := c.ShouldBindJSON(&function); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	if err := s.DB.Create(&function).Error; err != nil {
		s.internalError(c, "创建失败")
		return
	}

	c.JSON(http.StatusCreated, function)
}

// 更新功能点
func (s *Service) UpdateFunction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	var function models.Function
	if err := s.DB.First(&function, id).Error; err != nil {
		s.notFound(c, "功能点不存在")
		return
	}

	if err := c.ShouldBindJSON(&function); err != nil {
		s.badRequest(c, err.Error())
		return
	}

	if err := s.DB.Save(&function).Error; err != nil {
		s.internalError(c, "更新失败")
		return
	}

	c.JSON(http.StatusOK, function)
}

// 删除功能点
func (s *Service) DeleteFunction(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	if err := s.DB.Delete(&models.Function{}, id).Error; err != nil {
		s.internalError(c, "删除失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// 按分类获取功能点
func (s *Service) GetFunctionsByCategory(c *gin.Context) {
	category := c.Query("category")
	var functions []models.Function

	if category != "" {
		s.DB.Where("category = ?", category).Find(&functions)
	} else {
		s.DB.Find(&functions)
	}

	c.JSON(http.StatusOK, functions)
}
