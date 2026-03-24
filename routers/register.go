package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/handlers"
	"gorm.io/gorm"
)

func RegisterRouters(r *gin.Engine, db *gorm.DB) {
	s := handlers.NewService(db)

	// 产品类型管理
	r.GET("/type", s.ListProductTypes)
	r.GET("/type/:id", s.GetProductType)
	r.POST("/type", s.AddProductType)
	r.PUT("/type/:id", s.UpdateProductType)
	r.DELETE("/type/:id", s.DeleteProductType)

	// 功能点管理
	r.GET("/function", s.ListFunctions)
	r.GET("/function/:id", s.GetFunction)
	r.POST("/function", s.AddFunction)
	r.PUT("/function/:id", s.UpdateFunction)
	r.DELETE("/function/:id", s.DeleteFunction)
	r.GET("/functions/by-category", s.GetFunctionsByCategory)

	// 产品管理
	r.GET("/product", s.ListProducts)
	r.GET("/product/:id", s.GetProduct)
	r.POST("/product", s.AddProduct)
	r.PUT("/product/:id", s.UpdateProduct)
	r.DELETE("/product/:id", s.DeleteProduct)
	r.POST("/products/batch", s.GetProductsByIDs) // 批量获取产品

	// 网络拓扑管理
	r.GET("/topo", s.ListTopos)
	r.GET("/topo/:id", s.GetTopo)
	r.POST("/topo", s.AddTopo)
	r.PUT("/topo/:id", s.UpdateTopo)
	r.DELETE("/topo/:id", s.DeleteTopo)
	r.POST("/topo/:id/copy", s.CopyTopo)                     // 复制拓扑
	r.GET("/topo/:id/visualization", s.GetTopoVisualization) // 获取可视化数据

	// 冗余和空缺分析
	r.POST("/analyze", s.AnalyzeByProductIDs)
	r.POST("/analyze/by-topo/:id", s.AnalyzeByTopoID)

	// 产品建议
	r.POST("/suggest/:id", s.GetProductSuggestions)
}
