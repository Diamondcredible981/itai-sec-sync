package database

import (
	"fmt"

	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"github.com/iMayday-Yee/XinchuangAnalyze/utils"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	// 初始化产品类型
	productTypes := []models.ProductType{
		{ID: 1, Name: "防火墙", Description: "提供包过滤、状态检测等功能", Icon: "firewall", Color: "#ff6b6b"},
		{ID: 2, Name: "入侵检测系统", Description: "检测网络中异常行为", Icon: "shield", Color: "#4ecdc4"},
		{ID: 3, Name: "入侵防御系统", Description: "阻断入侵行为", Icon: "shield-check", Color: "#45b7d1"},
		{ID: 4, Name: "虚拟专用网", Description: "提供加密隧道支持", Icon: "vpn", Color: "#96ceb4"},
		{ID: 5, Name: "Web应用防火墙", Description: "保护Web应用免受攻击", Icon: "web", Color: "#ffeaa7"},
		{ID: 6, Name: "国产操作系统", Description: "支持自主可控的国产操作系统", Icon: "os", Color: "#dda0dd"},
		{ID: 7, Name: "国产数据库", Description: "支持自主研发的数据库系统", Icon: "database", Color: "#98d8c8"},
		{ID: 8, Name: "终端安全", Description: "终端病毒查杀与防护", Icon: "endpoint", Color: "#f7dc6f"},
		{ID: 9, Name: "安全信息与事件管理", Description: "集中管理与分析安全事件", Icon: "siem", Color: "#bb8fce"},
		{ID: 10, Name: "路由器", Description: "提供网络路由功能", Icon: "router", Color: "#85c1e9"},
		{ID: 11, Name: "交换机", Description: "提供网络交换功能", Icon: "switch", Color: "#82e0aa"},
	}

	for _, pt := range productTypes {
		db.FirstOrCreate(&pt, models.ProductType{ID: pt.ID})
	}

	// 初始化功能点
	functions := []models.Function{
		{ID: 101, Name: "数据包过滤", Description: "对网络数据包进行过滤", Category: "防护类"},
		{ID: 102, Name: "深度包检测", Description: "深度分析数据包内容", Category: "检测类"},
		{ID: 103, Name: "状态检测", Description: "基于连接状态的检测", Category: "检测类"},
		{ID: 104, Name: "病毒查杀", Description: "检测和清除恶意软件", Category: "防护类"},
		{ID: 105, Name: "威胁情报分析", Description: "基于威胁情报的分析", Category: "分析类"},
		{ID: 106, Name: "应用层控制", Description: "应用层访问控制", Category: "控制类"},
		{ID: 107, Name: "DDoS防护", Description: "分布式拒绝服务攻击防护", Category: "防护类"},
		{ID: 108, Name: "SSL解密", Description: "SSL/TLS流量解密", Category: "分析类"},
		{ID: 109, Name: "VPN功能", Description: "虚拟专用网络功能", Category: "连接类"},
		{ID: 110, Name: "身份认证", Description: "用户身份验证", Category: "认证类"},
		{ID: 111, Name: "访问控制", Description: "资源访问权限控制", Category: "控制类"},
		{ID: 112, Name: "日志采集", Description: "安全日志收集", Category: "管理类"},
		{ID: 113, Name: "日志分析", Description: "安全日志分析", Category: "分析类"},
		{ID: 114, Name: "攻击检测", Description: "攻击行为检测", Category: "检测类"},
		{ID: 115, Name: "行为分析", Description: "用户行为分析", Category: "分析类"},
		{ID: 116, Name: "零信任支持", Description: "零信任架构支持", Category: "架构类"},
	}

	for _, f := range functions {
		db.FirstOrCreate(&f, models.Function{ID: f.ID})
	}

	// 初始化产品数据
	products := []models.Product{
		// 防火墙产品
		{
			ID:             1,
			TypeID:         1,
			Name:           "天擎防火墙",
			Brand:          "奇安信",
			FunctionIDsStr: utils.IntSliceToString([]int{101, 103, 107, 111}),
		},
		{
			ID:             2,
			TypeID:         1,
			Name:           "山石网科防火墙",
			Brand:          "山石网科",
			FunctionIDsStr: utils.IntSliceToString([]int{101, 103, 106, 108, 111}),
		},
		{
			ID:             3,
			TypeID:         1,
			Name:           "深信服NGFW",
			Brand:          "深信服",
			FunctionIDsStr: utils.IntSliceToString([]int{101, 102, 103, 106, 107, 108, 111, 116}),
		},

		// 入侵检测系统
		{
			ID:             4,
			TypeID:         2,
			Name:           "绿盟IDS",
			Brand:          "绿盟科技",
			FunctionIDsStr: utils.IntSliceToString([]int{102, 105, 114, 115}),
		},
		{
			ID:             5,
			TypeID:         2,
			Name:           "启明星辰IDS",
			Brand:          "启明星辰",
			FunctionIDsStr: utils.IntSliceToString([]int{102, 103, 105, 114, 115}),
		},

		// 入侵防御系统
		{
			ID:             6,
			TypeID:         3,
			Name:           "安恒信息IPS",
			Brand:          "安恒信息",
			FunctionIDsStr: utils.IntSliceToString([]int{102, 103, 107, 114, 115}),
		},
		{
			ID:             7,
			TypeID:         3,
			Name:           "奇安信IPS",
			Brand:          "奇安信",
			FunctionIDsStr: utils.IntSliceToString([]int{101, 102, 103, 107, 114}),
		},

		// VPN产品
		{
			ID:             8,
			TypeID:         4,
			Name:           "深信服VPN网关",
			Brand:          "深信服",
			FunctionIDsStr: utils.IntSliceToString([]int{108, 109, 110, 111, 116}),
		},
		{
			ID:             9,
			TypeID:         4,
			Name:           "华为VPN设备",
			Brand:          "华为",
			FunctionIDsStr: utils.IntSliceToString([]int{109, 110, 111, 112}),
		},

		// WAF产品
		{
			ID:             10,
			TypeID:         5,
			Name:           "绿盟WAF",
			Brand:          "绿盟科技",
			FunctionIDsStr: utils.IntSliceToString([]int{102, 106, 108, 114}),
		},
		{
			ID:             11,
			TypeID:         5,
			Name:           "长亭雷池WAF",
			Brand:          "长亭科技",
			FunctionIDsStr: utils.IntSliceToString([]int{102, 106, 108, 113, 114}),
		},

		// 国产操作系统
		{
			ID:             12,
			TypeID:         6,
			Name:           "统信UOS",
			Brand:          "统信软件",
			FunctionIDsStr: utils.IntSliceToString([]int{104, 110, 111, 112}),
		},
		{
			ID:             13,
			TypeID:         6,
			Name:           "麒麟操作系统",
			Brand:          "中标麒麟",
			FunctionIDsStr: utils.IntSliceToString([]int{104, 110, 111, 112}),
		},
		{
			ID:             14,
			TypeID:         6,
			Name:           "欧拉openEuler",
			Brand:          "华为",
			FunctionIDsStr: utils.IntSliceToString([]int{104, 110, 111, 112, 113}),
		},

		// 国产数据库
		{
			ID:             15,
			TypeID:         7,
			Name:           "达梦数据库",
			Brand:          "达梦数据库",
			FunctionIDsStr: utils.IntSliceToString([]int{110, 111, 112, 113}),
		},
		{
			ID:             16,
			TypeID:         7,
			Name:           "人大金仓KingbaseES",
			Brand:          "人大金仓",
			FunctionIDsStr: utils.IntSliceToString([]int{110, 111, 112, 113}),
		},
		{
			ID:             17,
			TypeID:         7,
			Name:           "南大通用GBase",
			Brand:          "南大通用",
			FunctionIDsStr: utils.IntSliceToString([]int{110, 111, 112, 113}),
		},

		// 终端安全
		{
			ID:             18,
			TypeID:         8,
			Name:           "奇安信终端安全",
			Brand:          "奇安信",
			FunctionIDsStr: utils.IntSliceToString([]int{104, 114, 115}),
		},
		{
			ID:             19,
			TypeID:         8,
			Name:           "360安全卫士企业版",
			Brand:          "360",
			FunctionIDsStr: utils.IntSliceToString([]int{104, 114, 115, 112}),
		},

		// SIEM产品
		{
			ID:             20,
			TypeID:         9,
			Name:           "启明星辰SIEM",
			Brand:          "启明星辰",
			FunctionIDsStr: utils.IntSliceToString([]int{105, 112, 113, 114, 115}),
		},
		{
			ID:             21,
			TypeID:         9,
			Name:           "安恒信息SIEM",
			Brand:          "安恒信息",
			FunctionIDsStr: utils.IntSliceToString([]int{105, 112, 113, 114, 115, 116}),
		},

		// 路由器
		{
			ID:             22,
			TypeID:         10,
			Name:           "华为路由器",
			Brand:          "华为",
			FunctionIDsStr: utils.IntSliceToString([]int{101, 111, 112}),
		},
		{
			ID:             23,
			TypeID:         10,
			Name:           "锐捷路由器",
			Brand:          "锐捷网络",
			FunctionIDsStr: utils.IntSliceToString([]int{101, 111, 112}),
		},

		// 交换机
		{
			ID:             24,
			TypeID:         11,
			Name:           "华为交换机",
			Brand:          "华为",
			FunctionIDsStr: utils.IntSliceToString([]int{101, 111, 112}),
		},
		{
			ID:             25,
			TypeID:         11,
			Name:           "锐捷交换机",
			Brand:          "锐捷网络",
			FunctionIDsStr: utils.IntSliceToString([]int{101, 111, 112}),
		},
	}

	for _, product := range products {
		db.FirstOrCreate(&product, models.Product{ID: product.ID})
	}

	// 初始化网络拓扑数据（新图模式）
	topologies := []models.NetworkTopo{
		{ID: 1, Name: "企业网络安全拓扑"},
		{ID: 2, Name: "政务云安全拓扑"},
		{ID: 3, Name: "金融行业安全拓扑"},
		{ID: 4, Name: "互联网企业安全拓扑"},
		{ID: 5, Name: "教育行业安全拓扑"},
		{ID: 6, Name: "医疗行业安全拓扑"},
		{ID: 7, Name: "制造业安全拓扑"},
		{ID: 8, Name: "小型企业安全拓扑"},
		{ID: 9, Name: "零信任架构拓扑"},
		{ID: 10, Name: "云原生安全拓扑"},
		{ID: 11, Name: "超级冗余拓扑"},
	}

	topoProducts := map[uint][]int{
		1:  {1, 4, 8, 12, 15, 18, 22, 24},
		2:  {2, 5, 6, 10, 13, 16, 20, 23, 25},
		3:  {3, 4, 6, 7, 11, 14, 17, 19, 21, 22, 24},
		4:  {1, 2, 5, 10, 11, 12, 15, 18, 20, 22, 25},
		5:  {2, 4, 8, 9, 13, 16, 19, 23, 24},
		6:  {1, 3, 6, 10, 12, 15, 18, 21, 22, 24},
		7:  {2, 5, 7, 8, 14, 17, 19, 23, 25},
		8:  {1, 4, 12, 18, 22},
		9:  {3, 6, 8, 9, 11, 14, 17, 19, 21},
		10: {2, 3, 5, 10, 11, 13, 16, 20, 23, 25},
		11: {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
	}

	productMap := make(map[int]models.Product)
	for _, p := range products {
		productMap[int(p.ID)] = p
	}

	for _, topo := range topologies {
		db.FirstOrCreate(&topo, models.NetworkTopo{ID: topo.ID})

		// 已有图结构则不覆盖，避免重启时破坏用户拓扑
		var nodeCount int64
		db.Model(&models.TopoNode{}).Where("topo_id = ?", topo.ID).Count(&nodeCount)
		if nodeCount > 0 {
			continue
		}

		nodes, edges := buildSeedGraph(topo.ID, topoProducts[topo.ID], productMap)
		if len(nodes) > 0 {
			db.Create(&nodes)
		}
		if len(edges) > 0 {
			db.Create(&edges)
		}
	}
}

func buildSeedGraph(topoID uint, productIDs []int, productMap map[int]models.Product) ([]models.TopoNode, []models.TopoEdge) {
	nodes := make([]models.TopoNode, 0, len(productIDs))
	edges := make([]models.TopoEdge, 0)

	for i, pid := range productIDs {
		key := buildSeedNodeKey(topoID, pid, i+1)
		pidU := uint(pid)

		name := "设备"
		vendor := ""
		if p, ok := productMap[pid]; ok {
			name = p.Name
			vendor = p.Brand
		}

		nodeType := "hardware"
		zone := "core"
		if i < 2 {
			zone = "edge"
		} else if i > len(productIDs)-3 {
			zone = "internal"
		}

		nodes = append(nodes, models.TopoNode{
			TopoID:      topoID,
			NodeKey:     key,
			NodeType:    nodeType,
			Name:        name,
			Vendor:      vendor,
			ProductID:   &pidU,
			Criticality: "normal",
			Zone:        zone,
			Layer:       i + 1,
		})

		if i > 0 {
			edges = append(edges, models.TopoEdge{
				TopoID:      topoID,
				FromNodeKey: nodes[i-1].NodeKey,
				ToNodeKey:   key,
				EdgeType:    "network",
				Direction:   "uni",
				Weight:      1,
				Risk:        20,
			})
		}

		// 每4个节点增加一条跨层边，形成更接近真实图结构
		if i >= 3 && i%4 == 0 {
			edges = append(edges, models.TopoEdge{
				TopoID:      topoID,
				FromNodeKey: nodes[i-3].NodeKey,
				ToNodeKey:   key,
				EdgeType:    "network",
				Direction:   "uni",
				Weight:      2,
				Risk:        35,
			})
		}
	}

	return nodes, edges
}

func buildSeedNodeKey(topoID uint, productID int, index int) string {
	return fmt.Sprintf("topo-%d-p%d-n%d", topoID, productID, index)
}
