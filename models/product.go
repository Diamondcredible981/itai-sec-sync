package models

type Product struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	TypeID         uint   `json:"type_id"`
	Name           string `json:"name"`
	Brand          string `json:"brand"`
	FunctionIDs    []int  `gorm:"-" json:"function_ids"` // 程序逻辑处理
	FunctionIDsStr string `gorm:"type:text" json:"-"`    // 数据库存储
}
