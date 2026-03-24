package models

type NetworkTopo struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `json:"name"`
	ProductIDs    []int     `gorm:"-" json:"product_ids"` // 程序逻辑处理
	ProductIDsStr string    `gorm:"type:text" json:"-"`   // 数据库存储
	Products      []Product `gorm:"-" json:"products"`    // 前端展示用
}
