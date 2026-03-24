package models

type ProductType struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`  // 图标标识
	Color       string `json:"color"` // 显示颜色
}
