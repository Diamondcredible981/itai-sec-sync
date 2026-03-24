package models

type TopoNode struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	TopoID      uint   `gorm:"index:idx_topo_node_key,unique" json:"topo_id"`
	NodeKey     string `gorm:"size:64;index:idx_topo_node_key,unique" json:"node_key"`
	NodeType    string `gorm:"size:32;index" json:"node_type"` // hardware/software/os/service
	Name        string `json:"name"`
	Vendor      string `json:"vendor"`
	ProductID   *uint  `gorm:"index" json:"product_id,omitempty"`
	Criticality string `gorm:"size:16" json:"criticality"`
	Zone        string `gorm:"size:32" json:"zone"`
	Layer       int    `gorm:"default:0;index" json:"layer"`
	AttrsJSON   string `gorm:"type:text" json:"attrs_json,omitempty"`
}
