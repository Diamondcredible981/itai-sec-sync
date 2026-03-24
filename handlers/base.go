package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"gorm.io/gorm"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   APIError `json:"error"`
}

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

func (s *Service) badRequest(c *gin.Context, message string) {
	s.respondError(c, 400, "BAD_REQUEST", message)
}

func (s *Service) notFound(c *gin.Context, message string) {
	s.respondError(c, 404, "NOT_FOUND", message)
}

func (s *Service) internalError(c *gin.Context, message string) {
	s.respondError(c, 500, "INTERNAL_ERROR", message)
}

func (s *Service) loadTopoGraph(topoID uint) ([]models.TopoNode, []models.TopoEdge, error) {
	var nodes []models.TopoNode
	if err := s.DB.Where("topo_id = ?", topoID).Order("layer ASC, id ASC").Find(&nodes).Error; err != nil {
		return nil, nil, err
	}

	var edges []models.TopoEdge
	if err := s.DB.Where("topo_id = ?", topoID).Find(&edges).Error; err != nil {
		return nil, nil, err
	}

	return nodes, edges, nil
}

func (s *Service) deriveProductIDsFromTopo(topoID uint) ([]int, error) {
	var nodes []models.TopoNode
	if err := s.DB.Where("topo_id = ?", topoID).Order("layer ASC, id ASC").Find(&nodes).Error; err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return []int{}, nil
	}

	result := make([]int, 0)
	seen := make(map[int]bool)
	for _, n := range nodes {
		if n.ProductID == nil {
			continue
		}
		pid := int(*n.ProductID)
		if !seen[pid] {
			result = append(result, pid)
			seen[pid] = true
		}
	}

	return result, nil
}
