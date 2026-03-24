package models

type NetworkTopo struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	ProductIDs []int      `gorm:"-" json:"product_ids"` // 由图节点推导
	Products   []Product  `gorm:"-" json:"products"`    // 前端展示用
	Nodes      []TopoNode `gorm:"-" json:"nodes,omitempty"`
	Edges      []TopoEdge `gorm:"-" json:"edges,omitempty"`
}
