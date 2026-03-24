package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iMayday-Yee/XinchuangAnalyze/models"
	"github.com/iMayday-Yee/XinchuangAnalyze/utils"
)

type AttackPathNode struct {
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	Brand     string `json:"brand"`
	Type      string `json:"type"`
	Layer     int    `json:"layer"`
}

type AttackPathEdge struct {
	FromProductID   int    `json:"from_product_id"`
	ToProductID     int    `json:"to_product_id"`
	AttackTechnique string `json:"attack_technique"`
	RiskScore       int    `json:"risk_score"`
	Reason          string `json:"reason"`
}

type AttackPathResponse struct {
	TopologyID   uint             `json:"topology_id"`
	TopologyName string           `json:"topology_name"`
	SourceProductID *int          `json:"source_product_id,omitempty"`
	TargetProductID *int          `json:"target_product_id,omitempty"`
	PathLength   int              `json:"path_length"`
	OverallRisk  int              `json:"overall_risk"`
	RiskLevel    string           `json:"risk_level"`
	Summary      AttackPathSummary `json:"summary"`
	Nodes        []AttackPathNode `json:"nodes"`
	Edges        []AttackPathEdge `json:"edges"`
	KeyJumps     []AttackPathEdge `json:"key_jumps"`
	Mitigations  []string         `json:"mitigations"`
}

type AttackPathSummary struct {
	AvgJumpRisk        float64 `json:"avg_jump_risk"`
	MaxJumpRisk        int     `json:"max_jump_risk"`
	HighRiskJumpCount  int     `json:"high_risk_jump_count"`
	WeakestNodeID      int     `json:"weakest_node_id"`
	WeakestNodeName    string  `json:"weakest_node_name"`
}

func (s *Service) GetAttackPathByTopoID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.badRequest(c, "无效的ID")
		return
	}

	var topology models.NetworkTopo
	if err := s.DB.First(&topology, id).Error; err != nil {
		s.notFound(c, "网络拓扑不存在")
		return
	}
	topology.ProductIDs = utils.StringToIntSlice(topology.ProductIDsStr)

	sourceProductID, hasSource, err := parseOptionalInt(c.Query("source_product_id"))
	if err != nil {
		s.badRequest(c, "source_product_id 参数无效")
		return
	}
	targetProductID, hasTarget, err := parseOptionalInt(c.Query("target_product_id"))
	if err != nil {
		s.badRequest(c, "target_product_id 参数无效")
		return
	}
	if hasSource != hasTarget {
		s.badRequest(c, "source_product_id 和 target_product_id 需同时提供")
		return
	}

	var products []models.Product
	s.DB.Where("id IN ?", topology.ProductIDs).Find(&products)

	productMap := make(map[int]models.Product)
	for _, p := range products {
		p.FunctionIDs = utils.StringToIntSlice(p.FunctionIDsStr)
		productMap[int(p.ID)] = p
	}

	var productTypes []models.ProductType
	s.DB.Find(&productTypes)
	typeMap := make(map[uint]models.ProductType)
	for _, t := range productTypes {
		typeMap[t.ID] = t
	}

	orderedProducts := make([]models.Product, 0, len(topology.ProductIDs))
	for _, productID := range topology.ProductIDs {
		if p, ok := productMap[productID]; ok {
			orderedProducts = append(orderedProducts, p)
		}
	}

	var sourcePtr *int
	var targetPtr *int
	if hasSource && hasTarget {
		startIndex, endIndex := -1, -1
		for i, p := range orderedProducts {
			if int(p.ID) == sourceProductID {
				startIndex = i
			}
			if int(p.ID) == targetProductID {
				endIndex = i
			}
		}
		if startIndex == -1 || endIndex == -1 {
			s.badRequest(c, "指定的 source/target 产品不在该拓扑中")
			return
		}
		if startIndex >= endIndex {
			s.badRequest(c, "source_product_id 必须位于 target_product_id 之前")
			return
		}

		orderedProducts = orderedProducts[startIndex : endIndex+1]
		sourcePtr = &sourceProductID
		targetPtr = &targetProductID
	}

	if len(orderedProducts) < 2 {
		c.JSON(http.StatusOK, AttackPathResponse{
			TopologyID:       topology.ID,
			TopologyName:     topology.Name,
			SourceProductID:  sourcePtr,
			TargetProductID:  targetPtr,
			PathLength:       len(orderedProducts),
			OverallRisk:      0,
			RiskLevel:        "low",
			Summary:          AttackPathSummary{},
			Nodes:            []AttackPathNode{},
			Edges:            []AttackPathEdge{},
			KeyJumps:         []AttackPathEdge{},
			Mitigations:      []string{"当前路径节点不足以形成横向攻击路径，建议持续监控并保持最小暴露面"},
		})
		return
	}

	nodes := make([]AttackPathNode, 0, len(orderedProducts))
	for i, p := range orderedProducts {
		pt := typeMap[p.TypeID]
		nodes = append(nodes, AttackPathNode{
			ProductID: int(p.ID),
			Name:      p.Name,
			Brand:     p.Brand,
			Type:      pt.Name,
			Layer:     i + 1,
		})
	}

	edges := make([]AttackPathEdge, 0)
	for i := 0; i < len(orderedProducts)-1; i++ {
		src := orderedProducts[i]
		dst := orderedProducts[i+1]

		risk, technique, reason := estimateEdgeRisk(src, dst, typeMap)
		edges = append(edges, AttackPathEdge{
			FromProductID:   int(src.ID),
			ToProductID:     int(dst.ID),
			AttackTechnique: technique,
			RiskScore:       risk,
			Reason:          reason,
		})
	}

	overallRisk := calcPathRisk(edges)
	keyJumps := pickKeyJumps(edges)
	mitigations := buildPathMitigations(orderedProducts, edges)
	summary := buildPathSummary(orderedProducts, edges)

	c.JSON(http.StatusOK, AttackPathResponse{
		TopologyID:      topology.ID,
		TopologyName:    topology.Name,
		SourceProductID: sourcePtr,
		TargetProductID: targetPtr,
		PathLength:      len(orderedProducts),
		OverallRisk:     overallRisk,
		RiskLevel:       riskLevelByScore(overallRisk),
		Summary:         summary,
		Nodes:           nodes,
		Edges:           edges,
		KeyJumps:        keyJumps,
		Mitigations:     mitigations,
	})
}

func parseOptionalInt(v string) (int, bool, error) {
	if v == "" {
		return 0, false, nil
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return 0, true, err
	}
	return parsed, true, nil
}

func estimateEdgeRisk(src, dst models.Product, typeMap map[uint]models.ProductType) (int, string, string) {
	srcType := typeMap[src.TypeID].Name
	dstType := typeMap[dst.TypeID].Name

	risk := 35
	technique := "lateral-movement"
	reason := "存在相邻设备访问面，可形成横向移动"

	if src.Brand != "" && src.Brand == dst.Brand {
		risk += 15
		reason = "同厂商相邻设备可能共享配置习惯，横向利用成本更低"
	}

	if srcType == dstType && srcType != "" {
		risk += 10
		technique = "same-tier-pivot"
		reason = "同类型设备间更易形成同层横向扩散"
	}

	if len(src.FunctionIDs) >= 5 || len(dst.FunctionIDs) >= 5 {
		risk += 15
		reason = "高功能密度节点一旦失陷，攻击收益更高"
	}

	if len(src.FunctionIDs) == 0 || len(dst.FunctionIDs) == 0 {
		risk += 10
		reason = "设备功能映射缺失，存在未知暴露面"
	}

	if risk > 100 {
		risk = 100
	}

	return risk, technique, reason
}

func calcPathRisk(edges []AttackPathEdge) int {
	if len(edges) == 0 {
		return 0
	}

	total := 0
	maxEdge := 0
	for _, edge := range edges {
		total += edge.RiskScore
		if edge.RiskScore > maxEdge {
			maxEdge = edge.RiskScore
		}
	}

	avg := total / len(edges)
	overall := int(float64(avg)*0.7 + float64(maxEdge)*0.3)
	if overall > 100 {
		return 100
	}
	return overall
}

func pickKeyJumps(edges []AttackPathEdge) []AttackPathEdge {
	if len(edges) == 0 {
		return []AttackPathEdge{}
	}

	sortedEdges := make([]AttackPathEdge, len(edges))
	copy(sortedEdges, edges)
	sort.Slice(sortedEdges, func(i, j int) bool {
		return sortedEdges[i].RiskScore > sortedEdges[j].RiskScore
	})

	limit := 3
	if len(sortedEdges) < limit {
		limit = len(sortedEdges)
	}
	return sortedEdges[:limit]
}

func buildPathMitigations(products []models.Product, edges []AttackPathEdge) []string {
	mitigations := []string{
		"对关键跳点之间实施最小权限访问控制与南北向/东西向微隔离策略",
		"在高风险跳点部署行为检测与关联告警，缩短横向移动发现时间",
	}

	if len(products) >= 2 {
		sameBrandPair := false
		for i := 0; i < len(products)-1; i++ {
			if products[i].Brand != "" && products[i].Brand == products[i+1].Brand {
				sameBrandPair = true
				break
			}
		}
		if sameBrandPair {
			mitigations = append(mitigations, "关键相邻节点优先引入异构厂商组合，降低同源漏洞联动风险")
		}
	}

	highRiskEdge := false
	for _, edge := range edges {
		if edge.RiskScore >= 75 {
			highRiskEdge = true
			break
		}
	}
	if highRiskEdge {
		mitigations = append(mitigations, "对高风险路径启用跳板审计与双因素认证，阻断凭证滥用扩散")
	}

	return mitigations
}

func buildPathSummary(products []models.Product, edges []AttackPathEdge) AttackPathSummary {
	if len(edges) == 0 {
		return AttackPathSummary{}
	}

	total := 0
	maxRisk := 0
	highRiskCount := 0
	nodeRisk := make(map[int]int)

	for _, edge := range edges {
		total += edge.RiskScore
		if edge.RiskScore > maxRisk {
			maxRisk = edge.RiskScore
		}
		if edge.RiskScore >= 70 {
			highRiskCount++
		}

		if edge.RiskScore > nodeRisk[edge.FromProductID] {
			nodeRisk[edge.FromProductID] = edge.RiskScore
		}
		if edge.RiskScore > nodeRisk[edge.ToProductID] {
			nodeRisk[edge.ToProductID] = edge.RiskScore
		}
	}

	weakestNodeID := 0
	weakestNodeName := ""
	weakestNodeRisk := -1
	for _, p := range products {
		risk := nodeRisk[int(p.ID)]
		if risk > weakestNodeRisk {
			weakestNodeRisk = risk
			weakestNodeID = int(p.ID)
			weakestNodeName = p.Name
		}
	}

	avg := float64(total) / float64(len(edges))
	return AttackPathSummary{
		AvgJumpRisk:       attackPathRound2(avg),
		MaxJumpRisk:       maxRisk,
		HighRiskJumpCount: highRiskCount,
		WeakestNodeID:     weakestNodeID,
		WeakestNodeName:   weakestNodeName,
	}
}

func attackPathRound2(v float64) float64 {
	return float64(int(v*100+0.5)) / 100
}

func riskLevelByScore(score int) string {
	switch {
	case score >= 75:
		return "critical"
	case score >= 50:
		return "high"
	case score >= 25:
		return "medium"
	default:
		return "low"
	}
}
