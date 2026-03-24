package models

type TopoEdge struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	TopoID      uint   `gorm:"index" json:"topo_id"`
	FromNodeKey string `gorm:"size:64;index" json:"from_node_key"`
	ToNodeKey   string `gorm:"size:64;index" json:"to_node_key"`
	EdgeType    string `gorm:"size:32;index" json:"edge_type"` // network/install/depend/trust
	Direction   string `gorm:"size:16" json:"direction"`       // uni/bi
	Weight      int    `gorm:"default:1" json:"weight"`
	Risk        int    `gorm:"default:0" json:"risk"`
	AttrsJSON   string `gorm:"type:text" json:"attrs_json,omitempty"`
}
